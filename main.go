package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"
)

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
		log.Fatalf("error output JSON for record %s: %v\n", r, err)
		return
	}
	fmt.Printf("initial record %s\n", string(pprint))
}
