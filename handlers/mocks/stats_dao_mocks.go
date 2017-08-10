package mocks

import (
	"github.com/Sovianum/myTgtTest/model"
	"time"
)

type SuccessStatsDaoMock struct{}

func (*SuccessStatsDaoMock) Save(stats model.Stats) error {
	return nil
}

func (*SuccessStatsDaoMock) GetStatsSlice([]time.Time, string, int) (model.StatsSlice, error) {
	var row = model.Row{Id: 10, Age: 10, Sex: "M", Count: 100}
	var item = model.StatsItem{
		Date: model.QuotedDate(time.Now()),
		Rows: []model.Row{row},
	}
	var slice = model.StatsSlice{Items: []model.StatsItem{item}}

	return slice, nil
}
