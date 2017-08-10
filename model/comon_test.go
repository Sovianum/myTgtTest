package model

import (
	"errors"
	"strings"
	"testing"
)

func TestGetReaderFunc_PresenceFail(t *testing.T) {
	var presenceChecker = func([]byte) error { return errors.New("Not found") }
	var validator = func(interface{}) error { return nil }

	var readerFunc = getReaderFunc(presenceChecker, validator)
	var dest interface{}

	var err = readerFunc(strings.NewReader("{}"), &dest)
	if err == nil {
		t.Error("Had to crash")
	}
	if err.Error() != "Not found" {
		t.Errorf("Wrong error expected %v got %v", "\"Not found\"", err.Error())
	}
}

func TestGetReaderFunc_ValidatorFail(t *testing.T) {
	var presenceChecker = func([]byte) error { return nil }
	var validator = func(interface{}) error { return errors.New("Invalid") }

	var readerFunc = getReaderFunc(presenceChecker, validator)
	var dest interface{}

	var err = readerFunc(strings.NewReader("{}"), &dest)
	if err == nil {
		t.Error("Had to crash")
	}
	if err.Error() != "Invalid" {
		t.Errorf("Wrong error expected %v got %v", "\"Invalid\"", err.Error())
	}
}

func TestGetReaderFunc_ParseFail(t *testing.T) {
	var presenceChecker = func([]byte) error { return nil }
	var validator = func(interface{}) error { return nil }

	var readerFunc = getReaderFunc(presenceChecker, validator)
	var dest interface{}

	var err = readerFunc(strings.NewReader("{"), &dest)
	if err == nil {
		t.Error("Had to crash")
	}
}

func TestGetReaderFunc_Success(t *testing.T) {
	var presenceChecker = func([]byte) error { return nil }
	var validator = func(interface{}) error { return nil }

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
	if err.Error() != "some" {
		t.Errorf("Wrong error expected %v got %v", "\"some\"", err.Error())
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
