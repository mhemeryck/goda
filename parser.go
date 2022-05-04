package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

// Prefix error represents the error where the CODA line has a wrong prefix
type PrefixError struct {}

func (e *PrefixError) Error() string {
	return fmt.Sprintf("Wrong prefix")
}

// parse implements the common parse functionality, using reflection
// The type and value passed in shall represent the reflected type / value of the struct we want to parse the data into
func parse(s string, t reflect.Type, v reflect.Value) error {
	// Loop through each of the struct fields
	for i := 0; i < t.NumField(); i++ {
		tt := t.Field(i)
		vv := v.Field(i)
		// Get the offset and length from the annotated struct fields
		offset, err := strconv.Atoi(tt.Tag.Get("offset"))
		if err != nil {
			return err
		}
		length, err := strconv.Atoi(tt.Tag.Get("length"))
		if err != nil {
			return err
		}
		// Switch to correct parser by using the field type
		switch vv.Type() {
		case reflect.TypeOf(int(0)):
			value, err := strconv.Atoi(s[offset : offset+length])
			if err != nil {
				return err
			}
			vv.SetInt(int64(value))
		case reflect.TypeOf(string("")):
			value := strings.TrimSpace(s[offset : offset+length])
			vv.SetString(value)
		case reflect.TypeOf(bool(false)):
			value := s[offset:offset+length] != ""
			vv.SetBool(value)
		case reflect.TypeOf(time.Time{}):
			value, err := time.Parse("020106", s[offset:offset+length])
			if err != nil {
				return err
			}
			vv.Set(reflect.ValueOf(value))
		case reflect.TypeOf(decimal.Decimal{}):
			value, err := strconv.Atoi(s[offset : offset+length])
			if err != nil {
				return err
			}
			// Shift decimal 3 places
			balance := decimal.New(int64(value), -3)
			vv.Set(reflect.ValueOf(balance))
		}
	}
	return nil
}

// Parse is the generic record string line parser
// Each line passed in is a raw CODA record line
// The line is matched to its specific record struct and the data is read into it
func Parse(line string) (r Record, err error) {
	if strings.HasPrefix(line, "0") {
		r = &InitialRecord{}
	} else if strings.HasPrefix(line, "1") {
		r = &OldBalanceRecord{}
	} else if strings.HasPrefix(line, "21") {
		r = &MovementRecord1{}
	} else if strings.HasPrefix(line, "22") {
		r = &MovementRecord2{}
	} else if strings.HasPrefix(line, "23") {
		r = &MovementRecord3{}
	} else if strings.HasPrefix(line, "31") {
		r = &InformationRecord1{}
	} else if strings.HasPrefix(line, "32") {
		r = &InformationRecord2{}
	} else if strings.HasPrefix(line, "33") {
		r = &InformationRecord3{}
	} else if strings.HasPrefix(line, "8") {
		r = &NewBalanceRecord{}
	} else if strings.HasPrefix(line, "4") {
		r = &FreeCommunicationRecord{}
	} else if strings.HasPrefix(line, "9") {
		r = &TrailerRecord{}
	} else {
		return nil, &PrefixError{}
	}

	err = r.Parse(line)
	return r, err
}

type Record interface {
	Parse(string) error
}

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

// Parse populates an initial record from a string
func (r *InitialRecord) Parse(s string) error {
	if !strings.HasPrefix(s, "0") {
		return &PrefixError{}
	}

	err := parse(s, reflect.TypeOf(r).Elem(), reflect.ValueOf(r).Elem())
	if err != nil {
		return err
	}
	// Specific duplicate check
	r.IsDuplicate = s[16:17] == "D"
	return nil
}

// Parse populates an old balance record from a string
func (r *OldBalanceRecord) Parse(s string) error {
	if !strings.HasPrefix(s, "1") {
		return &PrefixError{}
	}

	return parse(s, reflect.TypeOf(r).Elem(), reflect.ValueOf(r).Elem())
}

// Parse populates a movement record from a string
func (r *MovementRecord1) Parse(s string) error {
	if !strings.HasPrefix(s, "21") {
		return &PrefixError{}
	}

	return parse(s, reflect.TypeOf(r).Elem(), reflect.ValueOf(r).Elem())
}

// Parse populates a movement record 2 from a string
func (r *MovementRecord2) Parse(s string) error {
	if !strings.HasPrefix(s, "22") {
		return &PrefixError{}
	}

	return parse(s, reflect.TypeOf(r).Elem(), reflect.ValueOf(r).Elem())
}

// Parse populates a movement record 3 from a string
func (r *MovementRecord3) Parse(s string) error {
	if !strings.HasPrefix(s, "23") {
		return &PrefixError{}
	}

	return parse(s, reflect.TypeOf(r).Elem(), reflect.ValueOf(r).Elem())
}

func (r *InformationRecord1) Parse(s string) error {
	if !strings.HasPrefix(s, "31") {
		return &PrefixError{}
	}

	return parse(s, reflect.TypeOf(r).Elem(), reflect.ValueOf(r).Elem())
}

func (r *InformationRecord2) Parse(s string) error {
	if !strings.HasPrefix(s, "32") {
		return &PrefixError{}
	}

	return parse(s, reflect.TypeOf(r).Elem(), reflect.ValueOf(r).Elem())
}

func (r *InformationRecord3) Parse(s string) error {
	if !strings.HasPrefix(s, "33") {
		return &PrefixError{}
	}

	return parse(s, reflect.TypeOf(r).Elem(), reflect.ValueOf(r).Elem())
}

func (r *NewBalanceRecord) Parse(s string) error {
	if !strings.HasPrefix(s, "8") {
		return &PrefixError{}
	}

	return parse(s, reflect.TypeOf(r).Elem(), reflect.ValueOf(r).Elem())
}

func (r *FreeCommunicationRecord) Parse(s string) error {
	if !strings.HasPrefix(s, "4") {
		return &PrefixError{}
	}

	return parse(s, reflect.TypeOf(r).Elem(), reflect.ValueOf(r).Elem())
}

func (r *TrailerRecord) Parse(s string) error {
	if !strings.HasPrefix(s, "9") {
		return &PrefixError{}
	}

	return parse(s, reflect.TypeOf(r).Elem(), reflect.ValueOf(r).Elem())
}
