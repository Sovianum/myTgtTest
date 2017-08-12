package handlers

import (
	"database/sql"
	"github.com/Sovianum/myTgtTest/mylog"
	"testing"
)

func TestNewDBEnv(t *testing.T) {
	var testData = []struct {
		logger *mylog.Logger
		db     *sql.DB
	}{
		{nil, nil},
		{new(mylog.Logger), nil},
		{nil, new(sql.DB)},
		{new(mylog.Logger), new(sql.DB)},
	}

	for i, item := range testData {
		var env = NewDBEnv(item.db, item.logger)
		if env.Logger != item.logger {
			t.Errorf("Error on logger when i = %d", i)
		}
	}
}
