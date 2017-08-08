package mocks

import (
	"errors"
	"github.com/Sovianum/myTgtTest/model"
)

type FailUserDAOMock struct{}

func (*FailUserDAOMock) Save(model.Registration) error {
	return errors.New("")
}

func (*FailUserDAOMock) Exists(uint) bool {
	return true
}
