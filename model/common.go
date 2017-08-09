package model

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"strings"
)

type UnmarshalFunc func(reader io.Reader, dest interface{}) error

func GetUnmarshaller(
	presenceChecker func([]byte) error,
	validator func(interface{}) error,
) UnmarshalFunc {
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

func checkPresence(data []byte, fields []string, errMessages []string) error {
	var m = make(map[string]interface{})
	var err = json.Unmarshal(data, &m)

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
