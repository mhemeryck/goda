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

// type FieldType int
//
// const (
// 	Time FieldType = iota
// 	Number
// 	Text
// 	Bool
// )
//
// type fieldSpec struct {
// 	Start int
// 	End   int
// 	Name  string
// 	Type  FieldType
// }
//
// var initialRecordSpec = [10]fieldSpec{
// 	{Start: 5, End: 11, Name: "CreationDate", Type: Time},
// }

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
	ReferenceType             int
	Reference                 string
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
	r.CreationDate, err = time.Parse("020106", s[5:11])
	if err != nil {
		return err
	}
	// Bank identification number
	r.BankIdentificationNumber, err = strconv.Atoi(s[11:14])
	if err != nil {
		return err
	}
	// Duplicate check
	r.IsDuplicate = string(s[16]) == "D"
	// Reference
	r.Reference = strings.TrimSpace(s[24:34])
	// Addressee
	r.Addressee = strings.TrimSpace(s[34:60])
	// BIC
	r.BIC = strings.TrimSpace(s[60:71])
	// Account holder reference
	r.AccountHolderReference, err = strconv.Atoi(s[71:82])
	if err != nil {
		return err
	}
	// Transaction reference
	r.TransactionReference = strings.TrimSpace(s[88:104])
	// Related reference
	r.RelatedReference = strings.TrimSpace(s[104:120])
	// Version code
	r.VersionCode, err = strconv.Atoi(string(s[127]))
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
	r.AccountStructure, err = strconv.Atoi(string(s[1]))
	if err != nil {
		return err
	}
	// Sequence number
	r.SerialNumber, err = strconv.Atoi(s[2:5])
	if err != nil {
		return err
	}
	// Account numner
	r.AccountNumber = strings.TrimSpace(s[5:42])
	// Old balance sign. False is credit, true is debit
	r.BalanceSign = string(s[42]) == "1"
	// Old balance
	balance, err := strconv.Atoi(s[43:58])
	if err != nil {
		return err
	}
	// Shift decimal 3 places
	r.OldBalance = decimal.New(int64(balance), -3)
	// Old balance date
	r.BalanceDate, err = time.Parse("020106", s[58:64])
	if err != nil {
		return err
	}
	// Account holder name
	r.AccountHolderName = strings.TrimSpace(s[64:90])
	// Account description
	r.AccountDescription = strings.TrimSpace(s[90:125])
	// Sequence number
	r.BankStatementSerialNumber, err = strconv.Atoi(s[125:128])
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
	r.SerialNumber, err = strconv.Atoi(s[2:6])
	if err != nil {
		return err
	}
	// Detail number
	r.DetailNumber, err = strconv.Atoi(s[6:10])
	if err != nil {
		return err
	}
	// Bank reference number
	r.BankReferenceNumber = strings.TrimSpace(s[10:31])
	// Movement sign
	r.BalanceSign = string(s[31]) == "1"
	// Balance
	balance, err := strconv.Atoi(s[32:47])
	if err != nil {
		return err
	}
	// Shift decimal 3 places
	r.Balance = decimal.New(int64(balance), -3)
	// Value date
	r.BalanceDate, err = time.Parse("020106", s[47:53])
	if err != nil {
		return err
	}

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
