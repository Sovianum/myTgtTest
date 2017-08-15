package dao

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Sovianum/myTgtTest/model"
	"strings"
	"time"
)

const (
	saveStats           = `
	INSERT INTO Stats (userId, ts, action)
  	  SELECT $1, date_trunc('day', CAST($2 AS TIMESTAMP)), $3
	ON CONFLICT (userid, action, ts) DO UPDATE SET counter = Stats.counter + 1;
	`
	newGetStatsTemplate = `
	SELECT c.id id, c.age age, c.sex sex, s.counter cnt, s.ts ts FROM
	  Client c
	  JOIN Stats s ON c.id = s.userId
	WHERE s.ts IN ( %s ) AND s.action = $%d
	ORDER BY s.ts, cnt DESC, id
	LIMIT $%d;
	`
)

type dbStatsDAO struct {
	db *sql.DB
}

func NewDBStatsDAO(db *sql.DB) StatsDAO {
	var result = new(dbStatsDAO)
	result.db = db
	return result
}

func (statsDao *dbStatsDAO) Save(s model.Stats) error {
	var _, err = statsDao.db.Exec(saveStats, s.User, time.Time(s.Timestamp), s.Action)
	return err
}

type getStatsOutputRow struct {
	id    uint
	age   uint
	sex   string
	count uint
	ts    time.Time
}

func (statsDao *dbStatsDAO) GetStatsSlice(dates []time.Time, action string, limit int) (model.StatsSlice, error) {
	if !model.IsValidAction(action) {
		return model.NewStatsSlice(), errors.New(model.StatsInvalidAction)
	}

	if len(dates) == 0 {
		return model.NewStatsSlice(), nil
	}

	var query = getStatsSelectQuery(len(dates))
	var args = getStatsSelectArgs(dates, action, limit)
	var rows, dbErr = statsDao.db.Query(query, args...)
	if dbErr != nil {
		return model.StatsSlice{}, dbErr
	}

	var queryResult = make([]getStatsOutputRow, 0)
	var err error
	for rows.Next() {
		var row = getStatsOutputRow{}
		err = rows.Scan(&row.id, &row.age, &row.sex, &row.count, &row.ts)
		if err != nil {
			break
		}

		queryResult = append(queryResult, row)
	}

	if err == nil {
		err = rows.Err()
	}
	if err != nil {
		return model.StatsSlice{}, err
	}

	var result = processGetStatsOutputRows(queryResult)

	return result, nil
}

// function assumes rows to contain valid values
func processGetStatsOutputRows(rows []getStatsOutputRow) model.StatsSlice {
	if len(rows) == 0 {
		return model.NewStatsSlice()
	}

	var result = model.NewStatsSlice()
	var currItem = model.NewItem(rows[0].ts)
	for _, row := range rows {
		if row.ts != time.Time(currItem.Date) {
			result.Items = append(result.Items, currItem)
			currItem = model.NewItem(row.ts)
		}

		currItem.Rows = append(currItem.Rows, model.Row{
			Id:    row.id,
			Age:   row.age,
			Sex:   row.sex,
			Count: row.count,
		})
	}

	if len(result.Items) == 0 || currItem.Date != result.Items[len(result.Items)-1].Date {
		result.Items = append(result.Items, currItem)
	}

	return result
}

func getStatsSelectArgs(dates []time.Time, action string, limit int) []interface{} {
	var result = make([]interface{}, 0)

	for i := range dates {
		result = append(result, dates[i])
	}
	result = append(result, action, limit)
	return result
}

// function assumes dateCount > 0
func getStatsSelectQuery(dateCount int) string {
	if dateCount <= 0 {
		panic(fmt.Sprintf("dateCount must not be less then 1 (got %d)", dateCount))
	}

	var dateNumSlice = make([]string, 0)

	for i := 1; i != dateCount+1; i++ {
		dateNumSlice = append(dateNumSlice, fmt.Sprintf("$%d", i))
	}

	var dateStr = strings.Join(dateNumSlice, ", ")
	return fmt.Sprintf(newGetStatsTemplate, dateStr, dateCount+1, dateCount+2)
}
