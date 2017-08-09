package model

import (
	"testing"
	"errors"
	"strings"
)

func TestGetReaderFunc_PresenceFail(t *testing.T) {
	var presenceChecker = func([]byte) error {return errors.New("Not found")}
	var validator =  func(interface{}) error {return nil}

	var readerFunc = getReaderFunc(presenceChecker, validator)
	var dest interface{}

	var err = readerFunc(strings.NewReader("{}"), &dest)
	if err == nil {
		t.Error("Had to crash")
	}
}

func TestGetReaderFunc_ValidatorFail(t *testing.T) {
	var presenceChecker = func([]byte) error {return nil}
	var validator =  func(interface{}) error {return errors.New("Invalid")}

	var readerFunc = getReaderFunc(presenceChecker, validator)
	var dest interface{}

	var err = readerFunc(strings.NewReader("{}"), &dest)
	if err == nil {
		t.Error("Had to crash")
	}
}

func TestGetReaderFunc_ParseFail(t *testing.T) {
	var presenceChecker = func([]byte) error {return nil}
	var validator =  func(interface{}) error {return nil}

	var readerFunc = getReaderFunc(presenceChecker, validator)
	var dest interface{}

	var err = readerFunc(strings.NewReader("{"), &dest)
	if err == nil {
		t.Error("Had to crash")
	}
}

func TestGetReaderFunc_Success(t *testing.T) {
	var presenceChecker = func([]byte) error {return nil}
	var validator =  func(interface{}) error {return nil}

	var readerFunc = getReaderFunc(presenceChecker, validator)
	var dest interface{}

	var err = readerFunc(strings.NewReader("{}"), &dest)
	if err != nil {
		t.Error(err)
	}
}

func TestCheckPresence_ParseFail(t *testing.T) {
	var jsonStr = "{\"some\":\"foo\", \"any\":\"bar\""
	var err = checkPresence([]byte(jsonStr), []string{"some"}, []string{"some"})
	if err == nil {
		t.Error("Had to crash")
	}
}

func TestCheckPresence_PresenceFail(t *testing.T) {
	var jsonStr = "{\"some\":\"foo\", \"any\":\"bar\"}"
	var err = checkPresence([]byte(jsonStr), []string{"som"}, []string{"some"})
	if err == nil {
		t.Error("Had to crash")
	}
}

func TestCheckPresence_Success(t *testing.T) {
	var jsonStr = "{\"some\":\"foo\", \"any\":\"bar\"}"
	var err = checkPresence([]byte(jsonStr), []string{"some"}, []string{"some"})
	if err != nil {
		t.Error(err)
	}
}

func TestCheckPresence_UnequalLength(t *testing.T) {
	var jsonStr = "{\"some\":\"foo\", \"any\":\"bar\"}"
	var err = checkPresence([]byte(jsonStr), []string{"some"}, []string{})
	if err == nil {
		t.Error("Had to crash")
	}
}
