package model

import (
	"fmt"
	"time"
)

type QuotedDate time.Time

func (q QuotedDate) MarshalJSON() ([]byte, error) {
	ts := time.Time(q).Format("2006-01-02")
	stamp := fmt.Sprintf("\"%s\"", ts)

	return []byte(stamp), nil
}

func (q *QuotedDate) UnmarshalJSON(b []byte) error {
	var inputS = string(b)
	var ts, err = time.Parse("\"2006-01-02\"", inputS)

	if err != nil {
		return err
	}

	*q = QuotedDate(ts)
	return nil
}

type StatsSlice struct {
	Items []StatsItem `json:"items"`
}

func NewStatsSlice() StatsSlice {
	return StatsSlice{Items: make([]StatsItem, 0)}
}

type StatsItem struct {
	Date QuotedDate `json:"date"`
	Rows []Row      `json:"rows"`
}

func NewItem(date time.Time) StatsItem {
	return StatsItem{Date: QuotedDate(date), Rows: make([]Row, 0)}
}

type Row struct {
	Id    uint   `json:"id"`
	Age   uint   `json:"age"`
	Sex   string `json:"sex"`
	Count uint   `json:"count"`
}
