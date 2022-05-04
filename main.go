package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

)

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

	for _, r := range records {
		// Pretty print
		pprint, err := json.MarshalIndent(r, "", "    ")
		if err != nil {
			log.Fatalf("Error output JSON for record %v: %v\n", r, err)
			return
		}
		fmt.Printf("record %T %s\n", r, string(pprint))
	}
}
