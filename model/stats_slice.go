package model

import (
	"time"
	"fmt"
)

type Calendar time.Time

func (c Calendar) MarshalJSON() ([]byte, error) {
	ts := time.Time(c).Format("2006-01-02")
	stamp := fmt.Sprintf("\"%s\"", ts)

	return []byte(stamp), nil
}

func (c *Calendar) UnmarshalJSON(b []byte) error {
	var layout = "2006-01-02"

	var inputS = string(b)
	var ts, err = time.Parse(layout, inputS[1:len(inputS)-1]) // slicing removes quotes

	if err != nil {
		return err
	}

	*c = Calendar(ts)
	return nil
}

type StatsSlice struct {
	Items []StatsItem `json:"items"`
}

func NewStatsSlice() StatsSlice {
	return StatsSlice{Items: make([]StatsItem, 0)}
}

type StatsItem struct {
	Date Calendar `json:"date"`
	Rows []Row    `json:"rows"`
}

func NewItem(date time.Time) StatsItem {
	return StatsItem{Date: Calendar(date), Rows: make([]Row, 0)}
}

type Row struct {
	Id    uint   `json:"id"`
	Age   uint   `json:"age"`
	Sex   string `json:"sex"`
	Count uint   `json:"count"`
}

//func (r *Row) UnmarshalJSON(b []byte) error {
//	type RowAlias Row
//	var ra RowAlias
//
//	var parseErr = json.Unmarshal(b, ra)
//	if parseErr != nil {
//		return parseErr
//	}
//
//	var intSex, _ = strconv.Atoi(ra.Sex)
//	var stringSex, err = DecodeSex(intSex)
//	if err != nil {
//		return err
//	}
//
//	r.Sex = stringSex
//	r.Count = ra.Count
//	r.Age = ra.Age
//	r.Id = ra.Id
//
//	return nil
//}
