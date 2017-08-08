package model

import (
	"io"
	"errors"
)

const (
	Login   = "login"
	Like    = "like"
	Comment = "comments"
	Exit    = "exit"

	StatsRequiredTs     = "\"st field required\""
	StatsRequiredUser   = "\"user field required\""
	StatsRequiredAction = "\"action field required\""
	StatsInvalidAction = "\"invalid action: must be one of following values (login, like, comments, exit)\""
)

type Stats struct {
	Timestamp Timestamp `json:"ts"`
	User      uint      `json:"user"`
	Action    string    `json:"action"`
}

func (s *Stats) UnmarshalJSON(reader io.Reader) error {
	var presenceChecker = func(data []byte) error {
		return checkPresence(
			data,
			[]string{"ts", "user", "action"},
			[]string{StatsRequiredTs, StatsRequiredUser, StatsRequiredAction},
		)
	}

	var validator = func(val interface{}) error {
		var stats = val.(*Stats)

		var condition = stats.Action == Login || stats.Action == Like || stats.Action == Comment || stats.Action == Exit
		if !condition {
			return errors.New(StatsInvalidAction)
		}

		return nil
	}

	return GetUnmarshaller(presenceChecker, validator)(reader, s)
}