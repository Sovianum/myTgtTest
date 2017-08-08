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
