package dao

import (
	"database/sql"
	"github.com/Sovianum/myTgtTest/model"
	"time"
	"sort"
)

const (
	saveStats = `INSERT INTO Stats (userId, ts, action) VALUES ($1, $2, $3)`
	getStats  = `
	SELECT c.id id, c.age age, c.sex sex, count(*) cnt FROM
	  Client c JOIN Stats s ON c.id = s.userId
	WHERE s.ts >= $1 AND s.ts < $2 AND s.action = $3
	GROUP BY c.id
	ORDER BY cnt DESC
	LIMIT $4;
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
	var dbSlice, parseErr = s.DBSlice()
	if parseErr != nil {
		return parseErr
	}

	var _, err = statsDao.db.Exec(saveStats, dbSlice...)
	return err
}

func (statsDao *dbStatsDAO) Get(dates []time.Time, action string, limit int) (model.StatsSlice, error) {
	if len(dates) == 0 {
		return model.NewStatsSlice(), nil
	}

	var datesCopy = make([]time.Time, len(dates))
	copy(datesCopy, dates)

	sort.Slice(datesCopy, func(i, j int) bool {
		return datesCopy[i].Before(datesCopy[j])
	})

	var err error
	var item model.StatsItem
	var result = model.NewStatsSlice()

	for _, date := range datesCopy {
		item, err = statsDao.getItem(date, action, limit)

		if err != nil {
			break
		}

		result.Items = append(result.Items, item)
	}

	return result, err
}

func (statsDao *dbStatsDAO) getItem(date time.Time, action string, limit int) (model.StatsItem, error) {
	var before = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	var after = before.Add(24 * time.Hour)

	var actionCode, codeErr = model.EncodeAction(action)
	if codeErr != nil {
		return model.StatsItem{}, codeErr
	}

	var rows, dbErr = statsDao.db.Query(getStats, before, after, actionCode, limit)
	if dbErr != nil {
		return model.StatsItem{}, dbErr
	}
	defer rows.Close()

	var result = model.NewItem(date)

	var err error
	for rows.Next() {
		var row = model.Row{}
		err = rows.Scan(&row.Id, &row.Age, &row.Sex, &row.Count)
		if err != nil {
			break
		}

		result.Rows = append(result.Rows, row)
	}

	if err == nil {
		err = rows.Err()
	}

	return result, err
}
