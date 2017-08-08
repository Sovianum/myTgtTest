package mocks

import "github.com/Sovianum/myTgtTest/model"

type NotExistUserDAOMock struct{}

func (*NotExistUserDAOMock) Save(model.Registration) error {
	return nil
}

func (*NotExistUserDAOMock) Exists(uint) bool {
	return false
}
