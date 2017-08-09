package mocks

import (
	"errors"
	"github.com/Sovianum/myTgtTest/model"
)

type ExistUserDAOMock struct{}

func (*ExistUserDAOMock) Save(model.Registration) error {
	return errors.New("")
}

func (*ExistUserDAOMock) Exists(uint) bool {
	return true
}

type NotExistUserDAOMock struct{}

func (*NotExistUserDAOMock) Save(model.Registration) error {
	return nil
}

func (*NotExistUserDAOMock) Exists(uint) bool {
	return false
}
