package model

import (
	"errors"
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

		if !IsValidAction(stats.Action) {
			return errors.New(StatsInvalidAction)
		}

		return nil
	}

	return GetUnmarshaller(presenceChecker, validator)(reader, s)
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

func DecodeAction(actionCode int) (string, error) {
	switch actionCode {
	case 0:
		return Login, nil
	case 1:
		return Like, nil
	case 2:
		return Comment, nil
	case 3:
		return Exit, nil
	default:
		return "", errors.New("Unknown action")
	}
}
