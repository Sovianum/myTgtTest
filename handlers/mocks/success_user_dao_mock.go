package mocks

import "github.com/Sovianum/myTgtTest/model"

type SuccessUserDAOMock struct{}

func (*SuccessUserDAOMock) Save(model.Registration) error {
	return nil
}

func (*SuccessUserDAOMock) Exists(uint) bool {
	return false
}
