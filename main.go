package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

// parseString is a helper function to clean up substrings
func parseString(s string) string {
	return strings.TrimSpace(s)
}

// parseDate is a helper function to get a date from a substring
// Date is typically in DD-MM-YY format
func parseDate(s string) (time.Time, error) {
	return time.Parse("020106", s)
}

// parseDecimal is a helper function to get a balance value in decimal format
// The string is typically 12 pos + 3
func parseDecimal(s string) (decimal.Decimal, error) {
	balance, err := strconv.Atoi(s)
	if err != nil {
		return decimal.Decimal{}, err
	}
	// Shift decimal 3 places
	return decimal.New(int64(balance), -3), nil
}

// parseInt is a helper function to capture a substring into an int
func parseInt(s string) (int, error) {
	return strconv.Atoi(s)
}

// parseSign is a helper function for getting the sign representation
// false / 0 is used for credit, true / 1 for debit
func parseSign(s string) bool {
	return s == "1"
}

// Record defines the generic interface for lines in a CODA file
type Record interface {
	Parse(string) error
}

// InitialRecord represents the first line of the CODA file
type InitialRecord struct {
	CreationDate             time.Time
	BankIdentificationNumber int
	IsDuplicate              bool
	Reference                string
	Addressee                string
	BIC                      string
	AccountHolderReference   int
	TransactionReference     string
	RelatedReference         string
	VersionCode              int
}

// OldBalanceRecord represents the old balance at the start of the CODA file
type OldBalanceRecord struct {
	AccountStructure          int
	SerialNumber              int
	AccountNumber             string
	BalanceSign               bool // True means debit / false credit
	OldBalance                decimal.Decimal
	BalanceDate               time.Time
	AccountHolderName         string
	AccountDescription        string
	BankStatementSerialNumber int
}

// TransactionRecord represents a single transaction in a CODA file
type TransactionRecord struct {
	SerialNumber              int
	DetailNumber              int
	BankReferenceNumber       string
	BalanceSign               bool // True means debit / false credit
	Balance                   decimal.Decimal
	BalanceDate               time.Time
	TransactionCode           int
	CommunicationType         int
	CommunicationZone         string
	BookingDate               time.Time
	BankStatementSerialNumber int
	GlobalisationCode         int
	TransactionSequence       bool
	InformationSequence       bool
}

// Parse parses a string into an InitialRecord
func (r *InitialRecord) Parse(s string) (err error) {
	// Check if it's an initial record
	if !strings.HasPrefix(s, "0") {
		return errors.New("Not an initial record")
	}

	// Creation date
	r.CreationDate, err = parseDate(s[5:11])
	if err != nil {
		return err
	}
	// Bank identification number
	r.BankIdentificationNumber, err = parseInt(s[11:14])
	if err != nil {
		return err
	}
	// Duplicate check
	r.IsDuplicate = s[16:17] == "D"
	// Reference
	r.Reference = parseString(s[24:34])
	// Addressee
	r.Addressee = parseString(s[34:60])
	// BIC
	r.BIC = parseString(s[60:71])
	// Account holder reference
	r.AccountHolderReference, err = parseInt(s[71:82])
	if err != nil {
		return err
	}
	// Transaction reference
	r.TransactionReference = parseString(s[88:104])
	// Related reference
	r.RelatedReference = parseString(s[104:120])
	// Version code
	r.VersionCode, err = parseInt(s[127:128])
	if err != nil {
		return err
	}
	return err
}

// Parse reads a string s into an OldBalanceRecord
func (r *OldBalanceRecord) Parse(s string) (err error) {
	// Check if it's an initial record
	if !strings.HasPrefix(s, "1") {
		return errors.New("Not an old balance record")
	}
	// Account structure
	r.AccountStructure, err = parseInt(s[1:2])
	if err != nil {
		return err
	}
	// Sequence number
	r.SerialNumber, err = parseInt(s[2:5])
	if err != nil {
		return err
	}
	// Account numner
	r.AccountNumber = parseString(s[5:42])
	// Old balance sign. False is credit, true is debit
	r.BalanceSign = parseSign(s[42:43])
	// Old balance
	r.OldBalance, err = parseDecimal(s[53:58])
	if err != nil {
		return err
	}
	// Old balance date
	r.BalanceDate, err = parseDate(s[58:64])
	if err != nil {
		return err
	}
	// Account holder name
	r.AccountHolderName = parseString(s[64:90])
	// Account description
	r.AccountDescription = parseString(s[90:125])
	// Sequence number
	r.BankStatementSerialNumber, err = parseInt(s[125:128])
	if err != nil {
		return err
	}
	return err
}

// Parse reads a string s into an OldBalanceRecord
func (r *TransactionRecord) Parse(s string) (err error) {
	// Check if it's an initial record
	if !strings.HasPrefix(s, "21") {
		return errors.New("Not a transaction record")
	}
	// Continuous sequence number
	r.SerialNumber, err = parseInt(s[2:6])
	if err != nil {
		return err
	}
	// Detail number
	r.DetailNumber, err = parseInt(s[6:10])
	if err != nil {
		return err
	}
	// Bank reference number
	r.BankReferenceNumber = parseString(s[10:31])
	// Movement sign
	r.BalanceSign = parseSign(s[31:32])
	// Balance
	r.Balance, err = parseDecimal(s[32:47])
	if err != nil {
		return err
	}
	// Value date
	r.BalanceDate, err = parseDate(s[47:53])
	if err != nil {
		return err
	}
	// Transaction code
	r.TransactionCode, err = parseInt(s[53:61])
	if err != nil {
		return err
	}
	// Communcation type: 0 none or unstructured / 1 structured
	r.CommunicationType, err = parseInt(s[61:62])
	if err != nil {
		return err
	}
	// Communication zone
	r.CommunicationZone = parseString(s[62:115])
	// Entry date
	r.BookingDate, err = parseDate(s[115:121])
	if err != nil {
		return err
	}
	// Sequence number
	r.BankStatementSerialNumber, err = parseInt(s[121:124])
	if err != nil {
		return err
	}
	// Globalisation code
	r.GlobalisationCode, err = parseInt(s[124:125])
	if err != nil {
		return err
	}
	// Next code: there is another transaction record
	r.TransactionSequence = parseString(s[125:126]) == "1"
	// There is another information record
	r.InformationSequence = parseString(s[127:128]) == "1"
	return err
}

func main() {
	records := []Record{}

	// Initial record
	sample := `0000002011830005        59501140  ACCOUNTANCY J DE KNIJF    BBRUBEBB   00412694022 00000                                       2`
	initialRecord := &InitialRecord{}
	err := initialRecord.Parse(sample)
	if err != nil {
		log.Fatalf("error parsing initial record: %s\n", err)
	}
	records = append(records, initialRecord)

	// Old balance record
	sample = `12001BE28310002350520                  EUR0000000001074020291217ACCOUNTANCY J DE KNIJF    Zichtrekening                      001`
	oldBalanceRecord := &OldBalanceRecord{}
	err = oldBalanceRecord.Parse(sample)
	if err != nil {
		log.Fatalf("error parsing old balance record: %s\n", err)
	}
	records = append(records, oldBalanceRecord)

	// Transaction record
	sample = `2100010000a29a791e4cff4000006  0000000001350000020118001500000Factuur nummer :20172720                             02011800101 0`
	transactionRecord := &TransactionRecord{}
	err = transactionRecord.Parse(sample)
	records = append(records, transactionRecord)

	// Pretty print
	for _, r := range records {
		pprint, err := json.MarshalIndent(r, "", "    ")
		if err != nil {
			log.Fatalf("Error output JSON for record %v: %v\n", r, err)
			return
		}
		fmt.Printf("record %T %s\n", r, string(pprint))
	}
}
