package main

import (
	"fmt"
	"reflect"
	"time"
)

type InitialRecord struct {
	CreationDate             time.Time
	BankIdentificationNumber int
	IsDuplicate              bool
	Reference                string
}

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

func main() {

	var x interface{}
	x = 3.14

	t := reflect.TypeOf(x)
	v := reflect.ValueOf(x)

	fmt.Printf("%v, %v\n", t, v)
}
