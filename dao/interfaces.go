package dao

import (
	"github.com/Sovianum/myTgtTest/model"
	"time"
)

type UserDAO interface {
	Save(r model.Registration) error
	Exists(id uint) bool
}

type StatsDAO interface {
	Save(s model.Stats) error
	Get(dates []time.Time, action string, limit int) (model.StatsSlice, error)
}
