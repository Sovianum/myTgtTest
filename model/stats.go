package model

import (
	"time"
	"io"
)

const (
	Login   = "login"
	Like    = "like"
	Comment = "comments"
	Exit    = "exit"
)

type Stats struct {
	Timestamp time.Time `json:"ts"`
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

	var validator = func(interface{}) error {
		return nil
	}

	return GetUnmarshaller(presenceChecker, validator)(reader, s)
}
