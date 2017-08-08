package dao

import (
	"errors"
	"github.com/Sovianum/myTgtTest/model"
)

const (
	UserAlreadyExists = "User already exists"
)

type MapUserDAO struct {
	m map[uint]model.Registration
}

func NewMapUserDAO() *MapUserDAO {
	var result = new(MapUserDAO)
	result.m = make(map[uint]model.Registration)
	return result
}

func (dao *MapUserDAO) Save(r model.Registration) error {
	var _, ok = dao.m[r.Id]
	if ok {
		return errors.New(UserAlreadyExists)
	}

	dao.m[r.Id] = r
	return nil
}

func (dao *MapUserDAO) Exists(id uint) bool {
	var _, ok = dao.m[id]
	return ok
}
