package model

import (
	"testing"
	"time"
)

func TestTimestamp_MarshalJSON_Success(t *testing.T) {
	var timeSample = time.Now()
	var correctStr = timeSample.Format("2006-01-02T15:04:05")

	var tmp = new(Timestamp)
	*tmp = Timestamp(timeSample)
	var gotBytes, _ = tmp.MarshalJSON()
	if correctStr != string(gotBytes) {
		t.Error("Parsed wrongly")
	}
}