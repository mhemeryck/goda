package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

// parse implements the common parse functionality, using reflection
// The type and value passed in shall represent the reflected type / value of the struct we want to parse the data into
func parse(s string, t reflect.Type, v reflect.Value) error {
	// fmt.Printf("Type %v value %v\n", t, v)
	// Loop through each of the struct fields
	for i := 0; i < t.NumField(); i++ {
		tt := t.Field(i)
		vv := v.Field(i)
		// fmt.Printf("Field %v: name %v, type %v tag %v value %v\n", i, tt.Name, tt.Type, tt.Tag, vv)
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
			// fmt.Printf("Value is %d\n", value)
			vv.SetInt(int64(value))
		case reflect.TypeOf(string("")):
			value := strings.TrimSpace(s[offset : offset+length])
			// fmt.Printf("Value is %s\n", value)
			vv.SetString(value)
		case reflect.TypeOf(bool(false)):
			value := s[offset:offset+length] != ""
			// fmt.Printf("Value is %t\n", value)
			vv.SetBool(value)
		case reflect.TypeOf(time.Time{}):
			value, err := time.Parse("020106", s[offset:offset+length])
			if err != nil {
				return err
			}
			// fmt.Printf("Value is %v\n", value)
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
// Each line passed is in is a raw CODA record line
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
	} else {
		return nil, nil
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

// Parse populates an initial record from a string
func (r *InitialRecord) Parse(s string) error {
	if !strings.HasPrefix(s, "0") {
		return errors.New("Wrong prefix")
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
		return errors.New("Wrong prefix")
	}

	return parse(s, reflect.TypeOf(r).Elem(), reflect.ValueOf(r).Elem())
}

// Parse populates a movement record from a string
func (r *MovementRecord1) Parse(s string) error {
	if !strings.HasPrefix(s, "21") {
		return errors.New("Wrong prefix")
	}

	return parse(s, reflect.TypeOf(r).Elem(), reflect.ValueOf(r).Elem())
}

// Parse populates a movement record 2 from a string
func (r *MovementRecord2) Parse(s string) error {
	if !strings.HasPrefix(s, "22") {
		return errors.New("Wrong prefix")
	}

	return parse(s, reflect.TypeOf(r).Elem(), reflect.ValueOf(r).Elem())
}

// Parse populates a movement record 3 from a string
func (r *MovementRecord3) Parse(s string) error {
	if !strings.HasPrefix(s, "23") {
		return errors.New("Wrong prefix")
	}

	return parse(s, reflect.TypeOf(r).Elem(), reflect.ValueOf(r).Elem())
}

//var sample = `0000002011830005        59501140  ACCOUNTANCY J DE KNIJF    BBRUBEBB   00412694022 00000                                       2`

const filename = "./sample.cod"

func main() {
	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		log.Fatalf("error opening file %s: %v\n", filename, err)
	}
	scanner := bufio.NewScanner(f)

	records := []Record{}
	var r Record
	for scanner.Scan() {
		line := scanner.Text()
		r, err = Parse(line)
		if err != nil {
			log.Fatalf("error parsing line %s: %v\n", line, err)
		}
		if r != nil {
			records = append(records, r)
		}
	}

	r = records[4]
	// Pretty print
	pprint, err := json.MarshalIndent(r, "", "    ")
	if err != nil {
		log.Fatalf("Error output JSON for record %v: %v\n", r, err)
		return
	}
	fmt.Printf("record %T %s\n", r, string(pprint))
}
