package model

import "time"

type StatsSlice struct {
	Items []StatsItem `json:"items"`
}

func NewStatsSlice() StatsSlice {
	return StatsSlice{Items:make([]StatsItem, 0)}
}

type StatsItem struct {
	Date Calendar `json:"date"`
	Rows []Row    `json:"rows"`
}

func NewItem(date time.Time) StatsItem {
	return StatsItem{Date:Calendar(date), Rows:make([]Row, 0)}
}

type Row struct {
	Id    uint   `json:"id"`
	Age   uint   `json:"age"`
	Sex   string `json:"sex"`
	Count uint   `json:"count"`
}
