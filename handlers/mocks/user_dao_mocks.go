package mocks

import (
	"errors"
	"github.com/Sovianum/myTgtTest/model"
)

type ExistUserDAOMock struct{}

func (*ExistUserDAOMock) Save(model.Registration) error {
	return errors.New("")
}

func (*ExistUserDAOMock) Exists(uint) (bool, error) {
	return true, nil
}

type NotExistUserDAOMock struct{}

func (*NotExistUserDAOMock) Save(model.Registration) error {
	return nil
}

func (*NotExistUserDAOMock) Exists(uint) (bool, error) {
	return false, nil
}
