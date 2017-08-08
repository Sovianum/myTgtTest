package model

import (
	"errors"
	"io"
	"strings"
)

const (
	MALE   = "M"
	FEMALE = "F"

	RegistrationRequiredId  = "\"id field required\""
	RegistrationRequiredAge = "\"age field required\""
	RegistrationRequiredSex = "\"sex field required\""
	RegistrationInvalidSex  = "\"invalid sex: must be either M or F\""
)

type Registration struct {
	Id  uint   `json:"id"`
	Age uint   `json:"age"`
	Sex string `json:"sex"`
}

func (r *Registration) UnmarshalJSON(reader io.Reader) error {
	var presenceChecker = func(data []byte) error {
		return checkPresence(
			data,
			[]string{"id", "age", "sex"},
			[]string{RegistrationRequiredId, RegistrationRequiredAge, RegistrationRequiredSex},
		)
	}

	var validator = func(val interface{}) error {
		var reg = val.(*Registration)
		return validateRegistration(reg)
	}

	return GetUnmarshaller(presenceChecker, validator)(reader, r)
}

func validateRegistration(r *Registration) error {
	var msgList = make([]string, 0)
	if r.Sex != MALE && r.Sex != FEMALE {
		msgList = append(msgList, RegistrationInvalidSex)
	}

	if len(msgList) != 0 {
		return errors.New(strings.Join(msgList, ";\n"))
	}
	return nil
}
