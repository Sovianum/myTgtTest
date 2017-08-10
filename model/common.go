package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

// This type represents a function, which reads data from reader to dest
type ReaderFunc func(reader io.Reader, dest interface{}) error

// This function returns ReaderFunc wrapped with two testing functions:
// 1. presenceChecker checks whether data read from reader contains all necessary fields
// 2. validator checks whether resulting object is in valid state after reading in data from reader
func getReaderFunc(
	presenceChecker func([]byte) error,
	validator func(interface{}) error,
) ReaderFunc {
	return func(reader io.Reader, dest interface{}) error {
		var data, readErr = ioutil.ReadAll(reader)
		if readErr != nil {
			return readErr
		}

		var absenceErr = presenceChecker(data)
		if absenceErr != nil {
			return absenceErr
		}

		var parseErr = json.Unmarshal(data, dest)
		if parseErr != nil {
			return parseErr
		}

		var validationErr = validator(dest)
		return validationErr
	}
}

// Function checks whether jsonData contains all fields from fields slice.
// errMessages slice contains messages which are used if some field is not found.
// Resulting error message consists of all corresponding errMessages, joined with ";\n"
func checkPresence(jsonData []byte, fields []string, errMessages []string) error {
	if len(fields) != len(errMessages) {
		return errors.New(
			fmt.Sprintf("Fields slice must have the same length (%v) as errMessages (%v)", len(fields), len(errMessages)),
		)
	}
	var m = make(map[string]interface{})
	var err = json.Unmarshal(jsonData, &m)

	if err != nil {
		return err
	}

	var messages = make([]string, 0)
	for i, field := range fields {
		_, ok := m[field]
		if !ok {
			messages = append(messages, errMessages[i])
		}
	}

	if len(messages) != 0 {
		return errors.New(strings.Join(messages, ";\n"))
	}
	return nil
}
