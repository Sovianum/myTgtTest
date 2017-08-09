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

type Timestamp time.Time

func (t *Timestamp) MarshalJSON() ([]byte, error) {
	ts := time.Time(*t).Format("2006-01-02T15:04:05")
	stamp := fmt.Sprint(ts)

	return []byte(stamp), nil
}

func (t *Timestamp) UnmarshalJSON(b []byte) error {
	var layout = "2006-01-02T15:04:05"

	var inputS = string(b)
	var ts, err = time.Parse(layout, inputS[1:len(inputS)-1]) // slicing removes quotes

	if err != nil {
		return err
	}

	*t = Timestamp(ts)
	return nil
}

type Stats struct {
	Timestamp Timestamp `json:"ts"`
	User      uint      `json:"user"`
	Action    string    `json:"action"`
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

	return GetReaderFunc(presenceChecker, validator)(reader, s)
}

func (s *Stats) DBSlice() ([]interface{}, error) {
	var encodedAction, err = EncodeAction(s.Action)
	if err != nil {
		return []interface{}{}, err
	}

	return []interface{}{
		s.User,
		time.Time(s.Timestamp),
		encodedAction,
	}, nil
}

func IsValidAction(action string) bool {
	return action == Login || action == Like || action == Comment || action == Exit
}

// Function encodes action string with int value to store it in database
func EncodeAction(action string) (int, error) {
	switch action {
	case Login:
		return 0, nil
	case Like:
		return 1, nil
	case Comment:
		return 2, nil
	case Exit:
		return 3, nil
	default:
		return -1, errors.New("Unknown action value")
	}
}

// TODO add type descriptions