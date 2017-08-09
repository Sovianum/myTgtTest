package dao

import (
	"database/sql"
	"github.com/Sovianum/myTgtTest/model"
)

const (
	saveUser = `INSERT INTO "User" (id, age, sex) VALUES ($1, $2, $3)`
	checkUser = `SELECT count(*) cnt FROM "User" u WHERE u.id = $1`
)

type dbUserDAO struct {
	db *sql.DB
}

func NewDBUserDAO(db *sql.DB) UserDAO {
	var result = new(dbUserDAO)
	result.db = db
	return result
}

func (dao *dbUserDAO) Save(r model.Registration) error {
	var args, sliceErr = r.DBSlice()
	if sliceErr != nil {
		return sliceErr
	}

	_, err := dao.db.Exec(saveUser, args...)
	return err
}

func (dao *dbUserDAO) Exists(id uint) bool {
	var cnt int
	var err = dao.db.QueryRow(checkUser, id).Scan(&cnt)
	if err != nil {
		panic(err)
	}

	return cnt > 0
}
