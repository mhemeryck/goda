package main

import (
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type GenericRecord interface{
	Parse(string) error
}

// Record is the common interface between all types of records (CODA lines)
type Record struct{}

// Parse is the generic record string parser
func (r *Record) Parse(s string) error {
	t := reflect.TypeOf(r).Elem()
	v := reflect.ValueOf(r).Elem()
	fmt.Printf("Type %v value %v\n", t, v)
	for i := 0; i < t.NumField(); i++ {
		tt := t.Field(i)
		vv := v.Field(i)
		fmt.Printf("Field %v: name %v, type %v tag %v value %v\n", i, tt.Name, tt.Type, tt.Tag, vv)
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
			fmt.Printf("Value is %d\n", value)
			vv.SetInt(int64(value))
		case reflect.TypeOf(string("")):
			value := strings.TrimSpace(s[offset : offset+length])
			fmt.Printf("Value is %s\n", value)
			vv.SetString(value)
		case reflect.TypeOf(bool(false)):
			value := s[offset:offset+length] != ""
			fmt.Printf("Value is %t\n", value)
			vv.SetBool(value)
		case reflect.TypeOf(time.Time{}):
			value, err := time.Parse("020106", s[offset:offset+length])
			if err != nil {
				return err
			}
			fmt.Printf("Value is %v\n", value)
			vv.Set(reflect.ValueOf(value))
		}
	}
	return nil
}

type InitialRecord struct {
	*Record
	CreationDate             time.Time `offset:"5" length:"6"`
	BankIdentificationNumber int       `offset:"11" length:"3"`
	IsDuplicate              bool      `offset:"16" length:"1"`
	Reference                string    `offset:"24" length:"10"`
}

//func (r *InitialRecord) Parse(s string) error {
//	return r.Parse(s)
//}

func parseDecimal(s string) (decimal.Decimal, error) {
	balance, err := strconv.Atoi(s)
	if err != nil {
		return decimal.Decimal{}, err
	}
	// Shift decimal 3 places
	return decimal.New(int64(balance), -3), nil
}

func printFields(x interface{}) {
	t := reflect.TypeOf(x)
	v := reflect.ValueOf(x)
	if t.Kind() != reflect.Struct {
		return
	}

	n := t.NumField()
	for i := 0; i < n; i++ {
		tt := t.Field(i)
		vv := v.Field(i)
		fmt.Printf("Field %v: name %v, type %v tag %v value %v\n", i, tt.Name, tt.Type, tt.Tag, vv)
	}
}

var sample = `0000002011830005        59501140  ACCOUNTANCY J DE KNIJF    BBRUBEBB   00412694022 00000                                       2`

func main() {
	// r := InitialRecord{time.Now(), 123, false, "456"}
	r := &InitialRecord{}
	err := r.Parse(sample)
	if err != nil {
		log.Fatalf("Could not parse sample %v: %v\n", r, err)
		return
	}

	// Pretty print
	pprint, err := json.MarshalIndent(r, "", "    ")
	if err != nil {
		log.Fatalf("Error output JSON for record %v: %v\n", r, err)
		return
	}
	fmt.Printf("record %T %s\n", r, string(pprint))
}
