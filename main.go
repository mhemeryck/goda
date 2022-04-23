package main

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type InitialRecord struct {
	CreationDate             time.Time `offset:"5" length:"6"`
	BankIdentificationNumber int       `offset:"11" length:"3"`
	IsDuplicate              bool      `offset:"16" length:"1"`
	Reference                string    `offset:"24" length:"10"`
}

var sample = `0000002011830005        59501140  ACCOUNTANCY J DE KNIJF    BBRUBEBB   00412694022 00000                                       2`

//for _, field := range reflect.VisibleFields(reflect.TypeOf(*transactionPurposeRecord)) {
//	if field.Type.Kind() == reflect.Int {
//		fmt.Printf("Found an int %s!\n", field.Name)
//	}
//	if field.Type.Kind() == reflect.ValueOf(time.Time{}).Kind() {
//		fmt.Printf("Found time %s!\n", field.Name)
//	}
//	positions := make(map[string]int)
//	tag := field.Tag.Get("coda")
//	if tag != "" {
//		fmt.Printf("%s: %s\n", field.Type, tag)
//		for _, descriptions := range strings.Split(tag, ",") {
//			s := strings.Split(descriptions, ":")
//			o, _ := strconv.Atoi(string(s[1]))
//			positions[s[0]] = o
//		}
//		fmt.Printf("%v\n", positions)
//	}
//	if field.Type.Kind() == reflect.Int {
//		fmt.Printf("Found an int %s!\n", field.Name)
//		field.SetInt(parseInt(sample[positions["offset"] : positions["offset"]+positions["length"]]))
//	}
//}

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

func main() {
	// r := InitialRecord{time.Now(), 123, false, "456"}
	r := InitialRecord{}
	// Get dynamic type and value, dereference pointer
	t := reflect.TypeOf(&r).Elem()
	v := reflect.ValueOf(&r).Elem()
	for i := 0; i < t.NumField(); i++ {
		tt := t.Field(i)
		vv := v.Field(i)
		offset, _ := strconv.Atoi(tt.Tag.Get("offset"))
		length, _ := strconv.Atoi(tt.Tag.Get("length"))
		s := sample[offset : offset+length]
		switch vv.Type() {
		case reflect.TypeOf(int(0)):
			value, err := strconv.Atoi(s)
			if err != nil {
				log.Fatalf("Error setting value\n")
			}
			fmt.Printf("Value is %d\n", value)
			vv.SetInt(int64(value))
		case reflect.TypeOf(string("")):
			value := strings.TrimSpace(s)
			fmt.Printf("Value is %s\n", value)
			vv.SetString(value)
		case reflect.TypeOf(bool(false)):
			value := s != ""
			fmt.Printf("Value is %t\n", value)
			vv.SetBool(value)
		case reflect.TypeOf(time.Time{}):
			value, err := time.Parse("020106", s)
			if err != nil {
				log.Fatalf("Error setting value\n")
			}
			fmt.Printf("Value is %v\n", value)
			vv.Set(reflect.ValueOf(value))
		}
	}
	//t := reflect.TypeOf(*r)
	//v := reflect.ValueOf(*r)
	//for _, field := range reflect.VisibleFields(t) {
	//	vv := v.FieldByIndex(field.Index)
	//	offset, _ := strconv.Atoi(field.Tag.Get("offset"))
	//	length, _ := strconv.Atoi(field.Tag.Get("length"))
	//	fmt.Printf("Field %s type %v value %v offset %d length %d\n", field.Name, field.Type, vv, offset, length)
	//	switch field.Type {
	//	case reflect.TypeOf(int(0)):
	//		value, err := strconv.Atoi(sample[offset : offset+length])
	//		if err != nil {
	//			log.Fatalf("Error setting value\n")
	//		}
	//		fmt.Printf("Value is %d\n", value)
	//		vv.SetInt(int64(value))
	//	case reflect.TypeOf(time.Time{}):
	//		fmt.Printf("FOUND A TIMESTAMP\n")
	//	}
	//}

	// Pretty print
	pprint, err := json.MarshalIndent(r, "", "    ")
	if err != nil {
		log.Fatalf("Error output JSON for record %v: %v\n", r, err)
		return
	}
	fmt.Printf("record %T %s\n", r, string(pprint))
}
