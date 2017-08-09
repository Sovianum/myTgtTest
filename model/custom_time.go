package model

import (
	"fmt"
	"time"
)

type Timestamp time.Time

func (t *Timestamp) MarshalJSON() ([]byte, error) {
	ts := time.Time(*t).Format("2006-01-02T15:04:05")
	stamp := fmt.Sprint(ts)

	return []byte(stamp), nil
}

func (t *Timestamp) UnmarshalJSON(b []byte) error {
	var layout = "2006-01-02T15:04:05"

	var inputS = string(b)
	var ts, err = time.Parse(layout, inputS[1:len(inputS)-1]) // slicing removes quotes TODO check if correct

	if err != nil {
		return err
	}

	*t = Timestamp(ts)
	return nil
}

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
