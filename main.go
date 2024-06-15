package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

var regularPattern = regexp.MustCompile(`([^:]+): ([^:]+)`)
var base64Pattern = regexp.MustCompile(`([^:]+):: ([^:]+)`)

type Record map[string][]string

func main() {
	var bytes []byte
	var err error
	if len(os.Args) < 2 {
		bytes, err = io.ReadAll(os.Stdin)
	} else {
		if os.Args[1] == "-h" || os.Args[1] == "--help" {
			fmt.Println("Parse output from ldapsearch from file or stdin")
			fmt.Printf("Usage: %s [<filename>]\n", os.Args[0])
			os.Exit(0)
		}
		filename := os.Args[1]
		bytes, err = os.ReadFile(filename)
		if err != nil {
			errorAndExit("Failed to read %s: %s", filename, err.Error())
		}
	}
	var records []Record
	record := make(Record)
	key := ""
	value := ""
	isBase64 := false
	for _, line := range strings.Split(string(bytes), "\n") {
		if strings.HasPrefix(line, " ") {
			// Continuation of previous line
			value += strings.TrimSpace(line)
			continue
		}
		if line == "" {
			// End of record
			insertIntoRecord(record, key, value, isBase64)
			if len(record) > 0 {
				records = append(records, record)
			}
			record = make(map[string][]string)
			key = ""
			value = ""
			continue
		}
		match := regularPattern.FindStringSubmatch(line)
		if match != nil {
			insertIntoRecord(record, key, value, isBase64)
			isBase64 = false
			key = match[1]
			value = match[2]
			continue
		}
		match = base64Pattern.FindStringSubmatch(line)
		if match != nil {
			insertIntoRecord(record, key, value, isBase64)
			isBase64 = true
			key = match[1]
			value = match[2]
			continue
		}
		errorAndExit("Unexpected line: %s", line)
		os.Exit(1)
	}
	data, err := json.Marshal(records)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to marshal records to json: %s\n", err.Error())
		os.Exit(1)
	}
	fmt.Println(string(data))
}

func errorAndExit(format string, args ...any) {
	format += "\n"
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

func insertIntoRecord(record Record, key, value string, isBase64 bool) {
	if key != "" && value != "" {
		if isBase64 {
			bytes, err := base64.StdEncoding.DecodeString(value)
			if err != nil {
				errorAndExit("Failed to decode base64 string %s: %s", value, err.Error())
			}
			value = string(bytes)
		}
		record[key] = append(record[key], value)
	}
}
