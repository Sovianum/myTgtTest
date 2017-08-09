package mocks

import (
	"github.com/Sovianum/myTgtTest/model"
	"time"
)

type SuccessStatsDaoMock struct {}

func (*SuccessStatsDaoMock) Save(stats model.Stats) error {
	return nil
}

func (*SuccessStatsDaoMock) Get(time.Time, string, uint) model.StatsSlice {
	var item = model.StatsItem{}
	item.Rows = make([]model.Row, 0)

	return model.StatsSlice{Items:[]model.StatsItem{item}}
}
