package dao

import (
	"database/sql"
	"github.com/Sovianum/myTgtTest/model"
	"sort"
	"time"
)

const (
	saveStats = `INSERT INTO Stats (userId, ts, action) SELECT $1, $2, code FROM Action WHERE str = $3`
	getStats  = `
	SELECT c.id id, c.age age, ss.str sex, count(*) cnt FROM
	  Client c
	  JOIN Sex ss ON ss.code = c.sex
	  JOIN Stats st ON c.id = st.userId
	  JOIN Action a ON a.code = st.action
	WHERE st.ts >= $1 AND st.ts < $2 AND a.str = $3
	GROUP BY c.id, ss.str
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
	var _, err = statsDao.db.Exec(saveStats, s.User, time.Time(s.Timestamp), s.Action)
	return err
}

func (statsDao *dbStatsDAO) GetStatsSlice(dates []time.Time, action string, limit int) (model.StatsSlice, error) {
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

	var rows, dbErr = statsDao.db.Query(getStats, before, after, action, limit)
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
