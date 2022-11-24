package goda

import (
	"testing"
	"time"
)

//func TestParseInitialRecord(t *testing.T) {
//	sample := `0000013020912605        YjeybrNhwgMichael Campbell          BBRUBEBB   03155032542                                             2`
//	r := &InitialRecord{}
//	err := r.Parse(sample)
//	if err != nil {
//		t.Fatalf("could not parse sample: %v", err)
//	}
//	if !r.CreationDate.Equal(time.Date(2009, 2, 13, 0, 0, 0, 0, time.UTC)) {
//		t.Fatalf("wrong date parsed")
//	}
//	if r.BankIdentificationNumber != 126 {
//		t.Fatalf("incorrect bank identification code parsed")
//	}
//	if r.Reference != "YjeybrNhwg" {
//		t.Fatalf("wrong reference parsed")
//	}
//	//if r.AccountHolderReference != 3155032542 {
//	//	t.Fatalf("issue parsing account holder reference")
//	//}
//	if r.BIC != "BBRUBEBB" {
//		t.Fatalf("issue parsing BIC code: got %s", r.BIC)
//	}
//	if r.IsDuplicate == true {
//		t.Fatalf("issue parsing duplicate code")
//	}
//	if r.Addressee != "Michael Campbell" {
//		t.Fatalf("issue parsing addressee")
//	}
//}

//func TestParseOldBalanceRecord(t *testing.T) {
//	sample := `10000                                     0000000550584847241114                                                             000`
//	r := &OldBalanceRecord{}
//	err := r.Parse(sample)
//	if err != nil {
//		t.Fatalf("could not parse sample: %v", err)
//	}
//}
//
//func TestParseTransactionRecord(t *testing.T) {
//	sample := `2139660000                     0000000160483785051100000000000                                                     21031600000 0`
//	r := &TransactionRecord{}
//	err := r.Parse(sample)
//	if err != nil {
//		t.Fatalf("could not parse sample: %v", err)
//	}
//}
//
//func TestParseTransactionPurposeRecord(t *testing.T) {
//	sample := `2268590000                                                                                                                   0 0`
//	r := &TransactionPurposeRecord{}
//	err := r.Parse(sample)
//	if err != nil {
//		t.Fatalf("could not parse sample: %v", err)
//	}
//
//}
//
//func TestParseTransactionDetailRecord(t *testing.T) {
//	sample := `2306190000                                     Mark Davis                                                                    0 0`
//
//	r := &TransactionDetailRecord{}
//	err := r.Parse(sample)
//	if err != nil {
//		t.Fatalf("could not parse sample: %v", err)
//	}
//	if r.AccountHolderName != "Mark Davis" {
//		t.Fatalf("expected account holder name to be Mark Davis, got %s", r.AccountHolderName)
//	}
//}
