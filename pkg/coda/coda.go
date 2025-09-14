package coda

import (
	// "io"
	"time"

	"github.com/shopspring/decimal"
)

const (
	HeaderIdentifier = 0
	recordLength     = 128
)

type LineWriter interface {
	WriteLine(line string) error
	Flush() error
}

type HeaderRecord struct {
	CreationDate             time.Time
	BankIdentificationNumber int // TODO: combines more information
	IsDuplicate              bool
	Reference                string
	Addressee                string
	BIC                      string
	IdentificationNumber     string
	TransactionReference     string
	RelatedReference         string
}

func MarshalHeaderRecord(h HeaderRecord, w LineWriter) error {
	w.WriteLine("%")
	return nil
}

type OldBalanceRecord struct {
	AccountStructure        int // TODO: combine with account number
	SequenceNUmberPaper     uint16
	AccountNumber           string
	OldBalance              decimal.Decimal
	OldBalanceDate          time.Time
	AccountHolderName       string
	AccountDescription      string
	SequenceNumberStatement int16
}

type Movement interface {
	MovementArticleCode() uint8
}

var _ Movement = (*MovementRecord1)(nil)
var _ Movement = (*MovementRecord2)(nil)
var _ Movement = (*MovementRecord3)(nil)

type MovementRecord1 struct{}

// MovementArticleCode implements Movement.
func (m *MovementRecord1) MovementArticleCode() uint8 {
	return 1
}

type MovementRecord2 struct{}

// MovementArticleCode implements Movement.
func (m *MovementRecord2) MovementArticleCode() uint8 {
	return 2
}

type MovementRecord3 struct{}

// MovementArticleCode implements Movement.
func (m *MovementRecord3) MovementArticleCode() uint8 {
	return 3
}

type Information interface {
	InformationArticleCode() uint8
}

var _ Information = (*InformationRecord1)(nil)
var _ Information = (*InformationRecord2)(nil)
var _ Information = (*InformationRecord3)(nil)

type InformationRecord1 struct{}

// InformationArticleCode implements Information.
func (m *InformationRecord1) InformationArticleCode() uint8 {
	return 1
}

type InformationRecord2 struct{}

// InformationArticleCode implements Information.
func (m *InformationRecord2) InformationArticleCode() uint8 {
	return 2
}

type InformationRecord3 struct{}

// InformationArticleCode implements Information.
func (m *InformationRecord3) InformationArticleCode() uint8 {
	return 3
}

type NewBalanceRecord struct{}
type FreeCommunicationRecord struct{}
type TrailerRecord struct{}

type Statement struct {
	Header            HeaderRecord
	OldBalance        OldBalanceRecord
	Movements         *[]Movement
	Informations      *[]Information
	NewBalance        *NewBalanceRecord
	FreeCommunication *FreeCommunicationRecord
	Trailer           TrailerRecord
}
