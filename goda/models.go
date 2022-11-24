package goda

import (
	"github.com/shopspring/decimal"
	"time"
)

type InitialRecord struct {
	CreationDate             time.Time `offset:"5" length:"6"`
	BankIdentificationNumber int       `offset:"11" length:"3"`
	IsDuplicate              bool      `offset:"16" length:"1"`
	Reference                string    `offset:"24" length:"10"`
	Addressee                string    `offset:"34" length:"26"`
	BIC                      string    `offset:"60" length:"11"`
	IdentificationNumber     int       `offset:"71" length:"11"`
	SeparateApplicationCode  int       `offset:"83" length:"5"`
	TransactionReference     string    `offset:"88" length:"16"`
	RelatedReference         string    `offset:"104" length:"16"`
	VersionCode              string    `offset:"127" length:"1"`
}

type OldBalanceRecord struct {
	AccountStructure        int             `offset:"1" length:"1"`
	SequenceNumberPaper     int             `offset:"2" length:"3"`
	AccountNumber           string          `offset:"5" length:"37"`
	OldBalanceSign          int             `offset:"42" length:"1"`
	OldBalance              decimal.Decimal `offset:"43" length:"15"`
	OldBalanceDate          time.Time       `offset:"58" length:"6"`
	AccountHolderName       string          `offset:"64" length:"26"`
	AccountDescription      string          `offset:"90" length:"35"`
	SequenceNumberStatement int             `offset:"125" length:"3"`
}

type MovementRecord1 struct {
	SequenceNumber      int             `offset:"2" length:"4"`
	DetailNumber        int             `offset:"6" length:"4"`
	BankReferenceNumber string          `offset:"10" length:"21"`
	MovementSign        int             `offset:"31" length:"1"`
	Amount              decimal.Decimal `offset:"32" length:"15"`
	ValueDate           time.Time       `offset:"47" length:"6"`
	TransactionCode     int             `offset:"53" length:"8"`
	CommunicationType   int             `offset:"61" length:"1"`
	CommunicationZone   string          `offset:"62" length:"53"`
	EntryDate           time.Time       `offset:"115" length:"6"`
	SequenceNumberPaper int             `offset:"121" length:"3"`
	GlobalisationCode   int             `offset:"124" length:"1"`
	NextCode            int             `offset:"125" length:"1"`
	LinkCode            int             `offset:"127" length:"1"`
}

type MovementRecord2 struct {
	SequenceNumber      int    `offset:"2" length:"4"`
	DetailNumber        int    `offset:"6" length:"4"`
	Communication       string `offset:"10" length:"53"`
	CustomerReference   string `offset:"63" length:"35"`
	CounterPartyBIC     string `offset:"98" length:"11"`
	RTransactionType    string `offset:"112" length:"1"`
	ISOReasonReturnCode string `offset:"113" length:"4"`
	CategoryPurpose     string `offset:"117" length:"4"`
	Purpose             string `offset:"121" length:"4"`
	NextCode            int    `offset:"125" length:"1"`
	LinkCode            int    `offset:"127" length:"1"`
}

type MovementRecord3 struct {
	SequenceNumber            int    `offset:"2" length:"4"`
	DetailNumber              int    `offset:"6" length:"4"`
	CounterPartyAccountNumber string `offset:"10" length:"37"`
	CounterPartyName          string `offset:"47" length:"35"`
	Communication             string `offset:"82" length:"43"`
	NextCode                  int    `offset:"125" length:"1"`
	LinkCode                  int    `offset:"127" length:"1"`
}

type InformationRecord1 struct {
	SequenceNumber         int    `offset:"2" length:"4"`
	DetailNumber           int    `offset:"6" length:"4"`
	BankReferenceNumber    string `offset:"10" length:"21"`
	TransactionCode        int    `offset:"31" length:"8"`
	CommunicationStructure int    `offset:"39" length:"1"`
	Communication          string `offset:"40" length:"73"`
	NextCode               int    `offset:"125" length:"1"`
	LinkCode               int    `offset:"127" length:"1"`
}

type InformationRecord2 struct {
	SequenceNumber int    `offset:"2" length:"4"`
	DetailNumber   int    `offset:"6" length:"4"`
	Communication  string `offset:"10" length:"105"`
	NextCode       int    `offset:"125" length:"1"`
	LinkCode       int    `offset:"127" length:"1"`
}

type InformationRecord3 struct {
	SequenceNumber int    `offset:"2" length:"4"`
	DetailNumber   int    `offset:"6" length:"4"`
	Communication  string `offset:"10" length:"90"`
	NextCode       int    `offset:"125" length:"1"`
	LinkCode       int    `offset:"127" length:"1"`
}

type NewBalanceRecord struct {
	SequenceNumber int             `offset:"1" length:"3"`
	AccountNumber  string          `offset:"4" length:"37"`
	NewBalanceSign int             `offset:"41" length:"1" `
	NewBalance     decimal.Decimal `offset:"42" length:"15"`
	NewBalanceDate time.Time       `offset:"57" length:"6"`
	LinkCode       int             `offset:"127" length:"1"`
}

type FreeCommunicationRecord struct {
	SequenceNumber    int    `offset:"2" length:"4"`
	DetailNumber      int    `offset:"6" length:"4"`
	FreeCommunication string `offset:"32" length:"80"`
	LinkCode          int    `offset:"127" length:"1"`
}

type TrailerRecord struct {
	NumberRecords    int             `offset:"16" length:"6"`
	DebitMovement    decimal.Decimal `offset:"22" length:"15"`
	CreditMovement   decimal.Decimal `offset:"37" length:"15"`
	MultipleFileCode int             `offset:"127" length:"1"`
}
