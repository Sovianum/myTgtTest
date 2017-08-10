package model

import (
	"errors"
	"fmt"
	"io"
	"time"
)

const (
	Login   = "login"
	Like    = "like"
	Comment = "comments"
	Exit    = "exit"

	StatsRequiredTs     = "\"st field required\""
	StatsRequiredUser   = "\"user field required\""
	StatsRequiredAction = "\"action field required\""
	StatsInvalidAction  = "\"invalid action: must be one of following values (login, like, comments, exit)\""
)

type QuotedTime time.Time

func (t *QuotedTime) MarshalJSON() ([]byte, error) {
	ts := time.Time(*t).Format("2006-01-02T15:04:05")
	stamp := fmt.Sprintf("\"%v\"", ts)

	return []byte(stamp), nil
}

func (t *QuotedTime) UnmarshalJSON(b []byte) error {
	var layout = "2006-01-02T15:04:05"

	var inputS = string(b)
	var ts, err = time.Parse(layout, inputS[1:len(inputS)-1]) // slicing removes quotes

	if err != nil {
		return err
	}

	*t = QuotedTime(ts)
	return nil
}

type Stats struct {
	Timestamp QuotedTime `json:"ts"`
	User      uint       `json:"user"`
	Action    string     `json:"action"`
}

func (s *Stats) ReadJsonIn(reader io.Reader) error {
	var presenceChecker = func(data []byte) error {
		return checkPresence(
			data,
			[]string{"ts", "user", "action"},
			[]string{StatsRequiredTs, StatsRequiredUser, StatsRequiredAction},
		)
	}

	var validator = func(val interface{}) error {
		var stats = val.(*Stats)

		if !IsValidAction(stats.Action) {
			return errors.New(StatsInvalidAction)
		}

		return nil
	}

	return getReaderFunc(presenceChecker, validator)(reader, s)
}

func IsValidAction(action string) bool {
	return action == Login || action == Like || action == Comment || action == Exit
}
