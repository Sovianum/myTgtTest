package dao

import (
	"errors"
	"fmt"
	"github.com/Sovianum/myTgtTest/model"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"strings"
	"testing"
	"time"
)

func TestDbStatsDAO_Save_Success(t *testing.T) {
	var db, mock, err = sqlmock.New()

	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	var s = model.Stats{}
	s.ReadJsonIn(strings.NewReader("{\"user\":1, \"action\":\"login\", \"ts\":\"2017-06-30T14:12:34\"}"))

	mock.
		ExpectExec("INSERT INTO").
		WithArgs(1, time.Time(s.Timestamp), s.Action).
		WillReturnResult(sqlmock.NewResult(1, 1))

	var statsDAO = NewDBStatsDAO(db)
	var saveErr = statsDAO.Save(s)

	if saveErr != nil {
		t.Error(saveErr.Error())
	}
}

func TestDbStatsDAO_Save_DBFail(t *testing.T) {
	var db, mock, err = sqlmock.New()

	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	var s = model.Stats{}
	s.ReadJsonIn(strings.NewReader("{\"user\":1, \"action\":\"login\", \"ts\":\"2017-06-30T14:12:34\"}"))

	mock.
		ExpectExec("INSERT INTO").
		WithArgs(1, time.Time(s.Timestamp), s.Action).
		WillReturnError(errors.New("Failed to save"))

	var statsDAO = NewDBStatsDAO(db)
	var saveErr = statsDAO.Save(s)

	if saveErr == nil {
		t.Error("Had to crash")
	}
	if saveErr.Error() != "Failed to save" {
		t.Errorf("Wrong error expected %v, got %v", "\"Failed to save\"", saveErr.Error())
	}
}

func TestDbStatsDAO_Get_IsSorted(t *testing.T) {
	var db, mock, err = sqlmock.New()

	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	var action = model.Like
	var limit = 10

	var date1 = time.Date(2003, 10, 17, 0, 0, 0, 0, time.UTC)
	var date2 = time.Date(2005, 10, 17, 0, 0, 0, 0, time.UTC)
	var rows = sqlmock.NewRows([]string{"id", "age", "sex", "cnt", "ts"}).
		AddRow(0, 10, "F", 100, date1).
		AddRow(1, 9, "M", 80, date2)

	mock.ExpectQuery("SELECT").
		WithArgs(date1, date2, action, limit).
		WillReturnRows(rows)

	var statsDAO = NewDBStatsDAO(db).(*dbStatsDAO)
	var statsSlice, sliceErr = statsDAO.GetStatsSlice(
		[]time.Time{date1, date2},
		model.Like, limit,
	)

	if sliceErr != nil {
		t.Error(sliceErr)
	}

	if len(statsSlice.Items) != 2 {
		t.Errorf("Incorrect row num: expected 2, got %v", len(statsSlice.Items))
	}
}

func TestDbStatsDAO_Get_Empty(t *testing.T) {
	var db, _, err = sqlmock.New()

	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	var limit = 10
	var statsDAO = NewDBStatsDAO(db).(*dbStatsDAO)
	var statsSlice, sliceErr = statsDAO.GetStatsSlice(
		[]time.Time{},
		model.Like, limit,
	)

	if sliceErr != nil {
		t.Error(sliceErr)
	}

	if len(statsSlice.Items) != 0 {
		t.Errorf("Incorrect row num: expected 0, got %v", len(statsSlice.Items))
	}
}

func TestDbStatsDAO_Get_DBFail(t *testing.T) {
	var db, mock, err = sqlmock.New()

	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	var action = model.Like
	var limit = 10

	var before = time.Date(2003, 10, 17, 0, 0, 0, 0, time.UTC)
	var after = before.Add(24 * time.Hour)

	mock.ExpectQuery("SELECT").
		WithArgs(before, after, action, limit).
		WillReturnError(errors.New("Failed to get"))

	var statsDAO = NewDBStatsDAO(db).(*dbStatsDAO)
	var _, sliceErr = statsDAO.GetStatsSlice(
		[]time.Time{before, after},
		model.Like, limit,
	)

	if sliceErr == nil {
		t.Error("Had to crash")
	}

	if sliceErr.Error() != "Failed to get" {
		t.Errorf("Wrong error expected %v got %v", "\"Failed to get\"", sliceErr.Error())
	}
}

func TestProcessGetStatsOutputRows_Empty(t *testing.T) {
	var testData = []getStatsOutputRow{}
	var output = processGetStatsOutputRows(testData)

	if len(output.Items) != 0 {
		t.Error(fmt.Sprintf("Expected %v, got %v", 0, len(output.Items)))
	}
}

func TestProcessGetStatsOutputRows_RepeatingData(t *testing.T) {
	var date1 = time.Now()
	var date2 = date1.Add(time.Hour)

	var testData = []getStatsOutputRow{
		{ts: date1, count: 10, sex: "F", age: 10, id: 100},
		{ts: date1, count: 10, sex: "F", age: 10, id: 100},
		{ts: date2, count: 10, sex: "F", age: 10, id: 100},
		{ts: date2, count: 10, sex: "F", age: 10, id: 100},
		{ts: date2, count: 10, sex: "F", age: 10, id: 100},
	}
	var output = processGetStatsOutputRows(testData)

	if len(output.Items) != 2 {
		t.Error(fmt.Sprintf("Expected %v, got %v", 2, len(output.Items)))
	}

	if len(output.Items[0].Rows) != 2 {
		t.Error(fmt.Sprintf("Expected %v, got %v", 2, len(output.Items[0].Rows)))
	}

	if len(output.Items[1].Rows) != 3 {
		t.Error(fmt.Sprintf("Expected %v, got %v", 3, len(output.Items[1].Rows)))
	}
}

func TestProcessGetStatsOutputRows_OneElement(t *testing.T) {
	var testData = []getStatsOutputRow{
		{ts: time.Now(), count: 10, sex: "F", age: 10, id: 100},
	}
	var output = processGetStatsOutputRows(testData)

	if len(output.Items) != 1 {
		t.Error(fmt.Sprintf("Expected %v, got %v", 1, len(output.Items)))
	}

	if time.Time(output.Items[0].Date) != testData[0].ts {
		t.Error(fmt.Sprintf("Expected %v, got %v", testData[0].ts, time.Time(output.Items[0].Date)))
	}

	if output.Items[0].Rows[0].Id != testData[0].id {
		t.Error(fmt.Sprintf("Expected %v, got %v", testData[0].id, output.Items[0].Rows[0].Id))
	}

	if output.Items[0].Rows[0].Count != testData[0].count {
		t.Error(fmt.Sprintf("Expected %v, got %v", testData[0].count, output.Items[0].Rows[0].Count))
	}

	if output.Items[0].Rows[0].Sex != testData[0].sex {
		t.Error(fmt.Sprintf("Expected %v, got %v", testData[0].sex, output.Items[0].Rows[0].Sex))
	}

	if output.Items[0].Rows[0].Age != testData[0].age {
		t.Error(fmt.Sprintf("Expected %v, got %v", testData[0].age, output.Items[0].Rows[0].Age))
	}
}

func TestGetStatsSelectArgs(t *testing.T) {
	var dates = []time.Time{time.Now(), time.Now().Add(time.Hour)}
	var action = "some"
	var limit = 100

	var args = getStatsSelectArgs(dates, action, limit)

	if len(args) != 4 {
		t.Error(fmt.Sprintf("Expected %v, got %v", 4, len(args)))
	}

	if gotTime := args[0].(time.Time); gotTime != dates[0] {
		t.Error(fmt.Sprintf("Expected %v, got %v", dates[0], gotTime))
	}

	if gotTime := args[1].(time.Time); gotTime != dates[1] {
		t.Error(fmt.Sprintf("Expected %v, got %v", dates[1], gotTime))
	}

	if gotAction := args[2].(string); gotAction != action {
		t.Error(fmt.Sprintf("Expected %v, got %v", action, gotAction))
	}

	if gotLimit := args[3].(int); gotLimit != limit {
		t.Error(fmt.Sprintf("Expected %v, got %v", action, gotLimit))
	}
}

func TestGetStatsSelectQuery(t *testing.T) {
	var correctQuery = `
	SELECT c.id id, c.age age, c.sex sex, s.counter cnt, s.ts ts FROM
	  Client c
	  JOIN Stats s ON c.id = s.userId
	WHERE s.ts IN ( $1, $2, $3 ) AND s.action = $4
	ORDER BY s.ts, cnt DESC, id
	LIMIT $5;
	`

	var gotQuery = getStatsSelectQuery(3)

	if correctQuery != gotQuery {
		t.Error(fmt.Sprintf("Expected \n \"%s\"\n got \n\"%s\"", correctQuery, gotQuery))
	}
}
