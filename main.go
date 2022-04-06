package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"
)

type FieldType int

const (
	Time FieldType = iota
	Number
	Text
	Bool
)

type fieldSpec struct {
	Start int
	End   int
	Name  string
	Type  FieldType
}

var initialRecordSpec = [10]fieldSpec{
	{Start: 5, End: 11, Name: "CreationDate", Type: Time},
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
	Free                     string
	TransactionReference     string
	VersionCode              int
}

// parseInitialRecord parses a string into an InitialRecord
func parseInitialRecord(r *InitialRecord, s string) (err error) {
	// Check if it's an initial record
	if string(s[0]) != "0" {
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
	r.Reference = s[24:34]
	return err
}

func main() {
	sample := `0000013020912605        YjeybrNhwgMichael Campbell          BBRUBEBB   03155032542                                             2`

	r := &InitialRecord{}
	err := parseInitialRecord(r, sample)
	if err != nil {
		log.Fatalf("error parsing initial record: %s\n", err)
	}
	pprint, err := json.MarshalIndent(r, "", "    ")
	if err != nil {
		log.Fatalf("Error output JSON for record %v: %v\n", r, err)
		return
	}
	fmt.Printf("initial record %s\n", string(pprint))
}
