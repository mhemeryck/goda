package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
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
	OldBalance                int
	BalanceDate               time.Time
	AccountHolderName         string
	AccountDescription        string
	BankStatementSerialNumber int
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

// parseOldBalanceRecord reads a string s into an OldBalanceRecord
func (r *OldBalanceRecord) Parse(s string) (err error) {
	// Check if it's an initial record
	if !strings.HasPrefix(s, "1") {
		return errors.New("Not an old balance record")
	}
	return err
}

func main() {
	records := []Record{}

	// Initial record
	sample := `0000013020912605        XXXXXXXXXXMichael Campbell          BBRUBEBB   03155032542                                             2`
	initialRecord := &InitialRecord{}
	err := initialRecord.Parse(sample)
	if err != nil {
		log.Fatalf("error parsing initial record: %s\n", err)
	}
	records = append(records, initialRecord)

	// Old balance record
	sample = `12001BE28310002350520                  EUR0000000001074020291217ACCOUNTANCY J DE KNIJF    Zichtrekening                      001`
	OldBalanceRecord := &OldBalanceRecord{}
	err = OldBalanceRecord.Parse(sample)
	if err != nil {
		log.Fatalf("error parsing old balance record: %s\n", err)
	}
	records = append(records, OldBalanceRecord)

	for _, r := range records {
		pprint, err := json.MarshalIndent(r, "", "    ")
		if err != nil {
			log.Fatalf("Error output JSON for record %v: %v\n", r, err)
			return
		}
		fmt.Printf("record %T %s\n", r, string(pprint))
	}
}
