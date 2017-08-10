package model

import (
	"strings"
	"testing"
	"time"
)

func TestQuotedTime_MarshalJSON_Success(t *testing.T) {
	var timeSample = time.Now()
	var correctStr = timeSample.Format("\"2006-01-02T15:04:05\"")

	var tmp = new(QuotedTime)
	*tmp = QuotedTime(timeSample)
	var gotBytes, _ = tmp.MarshalJSON()
	if correctStr != string(gotBytes) {
		t.Errorf("Parsed wrongly expected %v got %v", correctStr, string(gotBytes))
	}
}

func TestQuotedTime_UnmarshalJSON_Success(t *testing.T) {
	var correctTime = QuotedTime(time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC))
	var testTime = QuotedTime(time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC))
	var err = testTime.UnmarshalJSON([]byte("\"2006-01-02T15:04:05\""))

	if err != nil {
		t.Error(err)
	}

	if correctTime != testTime {
		t.Error("Unmarshaled wrongly")
	}
}

func TestQuotedTime_UnmarshalJSON_ParseError(t *testing.T) {
	var testTime = QuotedTime(time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC))
	var err = testTime.UnmarshalJSON([]byte("\"2006-01-67T15:04:05\""))

	if err == nil {
		t.Error("Had to crash")
	}
}

func TestStats_ReadJsonIn_ParseError(t *testing.T) {
	var reader = strings.NewReader("{\"ts\": \"2006-01-01T15:04:05\", \"user\": 1, action\": \"like\"")
	var s = Stats{}
	var err = s.ReadJsonIn(reader)
	if err == nil {
		t.Error("Had to crash")
	}
}

func TestStats_ReadJsonIn_IncompleteData(t *testing.T) {
	var reader = strings.NewReader("{\"ts\": \"2006-01-01T15:04:05\", \"action\": \"like\"}")
	var s = Stats{}
	var err = s.ReadJsonIn(reader)
	if err == nil {
		t.Error("Had to crash")
	}
	if err.Error() != StatsRequiredUser {
		t.Errorf("Wrong error: expected %v; got %v", StatsRequiredUser, err.Error())
	}
}

func TestStats_ReadJsonIn_InvalidStats(t *testing.T) {
	var reader = strings.NewReader("{\"ts\": \"2006-01-01T15:04:05\", \"user\": 1, \"action\": \"dfg\"}")
	var s = Stats{}
	var err = s.ReadJsonIn(reader)
	if err == nil {
		t.Error("Had to crash")
	}
	if err.Error() != StatsInvalidAction {
		t.Errorf("Wrong error: expected %v; got %v", StatsInvalidAction, err.Error())
	}
}

func TestStats_ReadJsonIn_Success(t *testing.T) {
	var reader = strings.NewReader("{\"ts\": \"2006-01-01T15:04:05\", \"user\": 1, \"action\": \"login\"}")
	var s = Stats{}
	var err = s.ReadJsonIn(reader)
	if err != nil {
		t.Error(err.Error())
	}

	if s.User != 1 || s.Action != Login || s.Timestamp != QuotedTime(time.Date(2006, 1, 1, 15, 4, 5, 0, time.UTC)) {
		t.Error("Wrong values")
	}
}
