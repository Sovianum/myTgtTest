package model

import (
	"testing"
	"time"
)

func TestQuotedDate_MarshalJSON_Success(t *testing.T) {
	var timeSample = time.Now()
	var correctStr = timeSample.Format("\"2006-01-02\"")

	var tmp = new(QuotedDate)
	*tmp = QuotedDate(timeSample)
	var gotBytes, _ = tmp.MarshalJSON()
	if correctStr != string(gotBytes) {
		t.Errorf("Parsed wrongly: expected %v, got %v", correctStr, string(gotBytes))
	}
}

func TestQuotedDate_UnmarshalJSON_Success(t *testing.T) {
	var correctDate = QuotedDate(time.Date(2006, 1, 2, 0, 0, 0, 0, time.UTC))
	var testDate = QuotedDate(time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC))
	var err = testDate.UnmarshalJSON([]byte("\"2006-01-02\""))

	if err != nil {
		t.Error(err)
	}

	if correctDate != testDate {
		t.Error("Unmarshaled badly")
	}
}

func TestQuotedDate_UnmarshalJSON_ParseError(t *testing.T) {
	var testDate = QuotedDate(time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC))
	var err = testDate.UnmarshalJSON([]byte("\"2006-01-67T15:04:05\""))

	if err == nil {
		t.Error("Had to crash")
	}
}

func TestNewItem(t *testing.T) {
	var item = NewItem(time.Now())
	if item.Rows == nil {
		t.Error("Item creation failed")
	}
}

func TestNewStatsSlice(t *testing.T) {
	var ss = NewStatsSlice()
	if ss.Items == nil {
		t.Error("Stats slice creation failed")
	}
}
