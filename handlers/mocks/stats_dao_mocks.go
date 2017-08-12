package mocks

import (
	"errors"
	"github.com/Sovianum/myTgtTest/model"
	"time"
)

type SuccessStatsDAOMock struct{}

func (*SuccessStatsDAOMock) Save(stats model.Stats) error {
	return nil
}

func (*SuccessStatsDAOMock) GetStatsSlice([]time.Time, string, int) (model.StatsSlice, error) {
	var row = model.Row{Id: 10, Age: 10, Sex: "M", Count: 100}
	var item = model.StatsItem{
		Date: model.QuotedDate(time.Now()),
		Rows: []model.Row{row},
	}
	var slice = model.StatsSlice{Items: []model.StatsItem{item}}

	return slice, nil
}

type FailStatsDAOMock struct{}

func (*FailStatsDAOMock) Save(stats model.Stats) error {
	return errors.New("Failed to save")
}

func (*FailStatsDAOMock) GetStatsSlice([]time.Time, string, int) (model.StatsSlice, error) {
	return model.StatsSlice{}, errors.New("Failed to select")
}
