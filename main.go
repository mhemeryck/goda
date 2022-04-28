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
func parse(s string, t reflect.Type, v reflect.Value) error {
	// fmt.Printf("Type %v value %v\n", t, v)
	for i := 0; i < t.NumField(); i++ {
		tt := t.Field(i)
		vv := v.Field(i)
		// fmt.Printf("Field %v: name %v, type %v tag %v value %v\n", i, tt.Name, tt.Type, tt.Tag, vv)
		offset, err := strconv.Atoi(tt.Tag.Get("offset"))
		if err != nil {
			return err
		}
		length, err := strconv.Atoi(tt.Tag.Get("length"))
		if err != nil {
			return err
		}
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
		}
	}
	return nil
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
}

type OldBalanceRecord struct {
	AccountStructure int    `offset:"1" length:"1"`
	SequenceNumber   int    `offset:"2" length:"3"`
	AccountNumber    string `offset:"5" length:"37"`
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

	err := parse(s, reflect.TypeOf(r).Elem(), reflect.ValueOf(r).Elem())
	if err != nil {
		return err
	}
	return nil
}

// Parse is the generic record string parser
func Parse(line string) (r Record, err error) {
	if strings.HasPrefix(line, "0") {
		r = &InitialRecord{}
	} else if strings.HasPrefix(line, "1") {
		r = &OldBalanceRecord{}
	} else {
		return nil, nil
	}

	err = r.Parse(line)
	return r, err
}

func parseDecimal(s string) (decimal.Decimal, error) {
	balance, err := strconv.Atoi(s)
	if err != nil {
		return decimal.Decimal{}, err
	}
	// Shift decimal 3 places
	return decimal.New(int64(balance), -3), nil
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

	r = records[1]
	// Pretty print
	pprint, err := json.MarshalIndent(r, "", "    ")
	if err != nil {
		log.Fatalf("Error output JSON for record %v: %v\n", r, err)
		return
	}
	fmt.Printf("record %T %s\n", r, string(pprint))
}
