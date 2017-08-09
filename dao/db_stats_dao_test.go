package dao

import (
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
		WithArgs(1, time.Time(s.Timestamp), 0).
		WillReturnResult(sqlmock.NewResult(1, 1))

	var statsDAO = NewDBStatsDAO(db)
	var saveErr = statsDAO.Save(s)

	if saveErr != nil {
		t.Error(saveErr.Error())
	}
}

// TODO think of fail cases

func TestDbStatsDAO_GetItem_Success(t *testing.T) {
	var db, mock, err = sqlmock.New()

	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	var rows = sqlmock.NewRows([]string{"id", "age", "sex", "cnt"}).
		AddRow(1, 100, 0, 10).
		AddRow(2, 200, 1, 7)

	var before = time.Date(2005, 10, 17, 0, 0, 0, 0, time.UTC)
	var after = before.Add(24 * time.Hour)
	var action, _ = model.EncodeAction(model.Like)
	var limit = 10

	mock.ExpectQuery("SELECT").
		WithArgs(before, after, action, limit).
		WillReturnRows(rows)

	var statsDAO = NewDBStatsDAO(db).(*dbStatsDAO)
	var item, itemErr = statsDAO.getItem(before, model.Like, limit)

	if itemErr != nil {
		t.Error(itemErr)
	}

	if len(item.Rows) != 2 {
		t.Errorf("Incorrect row num: expected 2, got %v", len(item.Rows))
	}
}

func TestDbStatsDAO_GetItem_Empty(t *testing.T) {
	var db, mock, err = sqlmock.New()

	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	var rows = sqlmock.NewRows([]string{"id", "age", "sex", "cnt"})

	var before = time.Date(2005, 10, 17, 0, 0, 0, 0, time.UTC)
	var after = before.Add(24 * time.Hour)
	var action, _ = model.EncodeAction(model.Like)
	var limit = 10

	mock.ExpectQuery("SELECT").
		WithArgs(before, after, action, limit).
		WillReturnRows(rows)

	var statsDAO = NewDBStatsDAO(db).(*dbStatsDAO)
	var item, itemErr = statsDAO.getItem(before, model.Like, limit)

	if itemErr != nil {
		t.Error(itemErr)
	}

	if len(item.Rows) != 0 {
		t.Errorf("Incorrect row num: expected 0, got %v", len(item.Rows))
	}
}

func TestDbStatsDAO_GetItem_BadAction(t *testing.T) {
	var db, mock, err = sqlmock.New()

	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	var rows = sqlmock.NewRows([]string{"id", "age", "sex", "cnt"})

	var before = time.Date(2005, 10, 17, 0, 0, 0, 0, time.UTC)
	var after = before.Add(24 * time.Hour)
	var action, _ = model.EncodeAction(model.Like)
	var limit = 10

	mock.ExpectQuery("SELECT").
		WithArgs(before, after, action, limit).
		WillReturnRows(rows)

	var statsDAO = NewDBStatsDAO(db).(*dbStatsDAO)
	var _, itemErr = statsDAO.getItem(before, "", limit)

	if itemErr == nil {
		t.Error("Had to crash on strange action")
	}
}

func TestDbStatsDAO_Get_IsSorted(t *testing.T) {
	var db, mock, err = sqlmock.New()

	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	var action, _ = model.EncodeAction(model.Like)
	var limit = 10

	var testData = []struct {
		before time.Time
		after  time.Time
		rows   *sqlmock.Rows
	}{
		{
			before: time.Date(2003, 10, 17, 0, 0, 0, 0, time.UTC),
			after:  time.Date(2003, 10, 18, 0, 0, 0, 0, time.UTC),
			rows: sqlmock.NewRows([]string{"id", "age", "sex", "cnt"}).
				AddRow(0, 10, 1, 100).
				AddRow(1, 9, 0, 80),
		},
		{
			before: time.Date(2005, 10, 17, 0, 0, 0, 0, time.UTC),
			after:  time.Date(2005, 10, 18, 0, 0, 0, 0, time.UTC),
			rows: sqlmock.NewRows([]string{"id", "age", "sex", "cnt"}).
				AddRow(0, 100, 0, 10).
				AddRow(1, 99, 1, 8),
		},
	}

	for _, dataItem := range testData {
		mock.ExpectQuery("SELECT").
			WithArgs(dataItem.before, dataItem.after, action, limit).
			WillReturnRows(dataItem.rows)
	}

	var statsDAO = NewDBStatsDAO(db).(*dbStatsDAO)
	var statsSlice, sliceErr = statsDAO.Get(
		[]time.Time{testData[1].before, testData[0].before},
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
	var statsSlice, sliceErr = statsDAO.Get(
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
