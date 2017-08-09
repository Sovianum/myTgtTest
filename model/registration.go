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

func (r *Registration) DBSlice() ([]interface{}, error) {
	var encodedSex, err = EncodeSex(r.Sex)
	if err != nil {
		return []interface{}{}, err
	}

	return []interface{}{
		r.Id,
		r.Age,
		encodedSex,
	}, nil
}

//func FromDBSlice(id uint, age uint, sexCode int) (Registration, error) {
//	var sex, err = DecodeSex(sexCode)
//	if err != nil {
//		return Registration{}, err
//	}
//
//	return Registration{Id:id, Age:age, Sex:sex}, nil
//}


func EncodeSex(sexString string) (int, error) {
	switch sexString {
	case MALE:
		return 0, nil
	case FEMALE:
		return 1, nil
	default:
		return -1, errors.New("Strange sex")
	}
}

func DecodeSex(sexNum int) (string, error) {
	switch sexNum {
	case 0:
		return MALE, nil
	case 1:
		return FEMALE, nil
	default:
		return "", errors.New("Strange value")
	}
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
