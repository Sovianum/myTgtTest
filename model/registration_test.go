package model

import (
	"strings"
	"testing"
)

func TestRegistration_ReadJsonIn_ParseError(t *testing.T) {
	var reader = strings.NewReader("{\"age\": 10, \"sex\": \"F\"")
	var r = Registration{}
	var err = r.ReadJsonIn(reader)
	if err == nil {
		t.Error("Had to crash")
	}
}

func TestRegistration_ReadJsonIn_IncompleteData(t *testing.T) {
	var reader = strings.NewReader("{\"age\": 10, \"sex\": \"F\"}")
	var r = Registration{}
	var err = r.ReadJsonIn(reader)
	if err == nil {
		t.Error("Had to crash")
	}
	if err.Error() != RegistrationRequiredId {
		t.Errorf("Wrong error: expected %v, got %v", RegistrationRequiredId, err.Error())
	}
}

func TestRegistration_ReadJsonIn_InvalidRegistration(t *testing.T) {
	var reader = strings.NewReader("{\"id\": 10, \"age\": 10, \"sex\": \"p\"}")
	var r = Registration{}
	var err = r.ReadJsonIn(reader)
	if err == nil {
		t.Error("Had to crash")
	}
	if err.Error() != RegistrationInvalidSex {
		t.Errorf("Wrong error: expected %v, got %v", RegistrationInvalidSex, err.Error())
	}
}

func TestRegistration_ReadJsonIn_Success(t *testing.T) {
	var reader = strings.NewReader("{\"id\": 10, \"age\": 10, \"sex\": \"F\"}")
	var r = Registration{}
	var err = r.ReadJsonIn(reader)
	if err != nil {
		t.Error(err)
	}
	if r.Sex != "F" || r.Age != 10 || r.Id != 10 {
		t.Error("Not set correct values")
	}
}
