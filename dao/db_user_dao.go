package dao

import (
	"database/sql"
	"github.com/Sovianum/myTgtTest/model"
)

const (
	saveUser  = `INSERT INTO Client (id, age, sex) SELECT $1, $2, code FROM Sex WHERE str = $3`
	checkUser = `SELECT count(*) cnt FROM Client u WHERE u.id = $1`
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
	_, err := dao.db.Exec(saveUser, r.Id, r.Age, r.Sex)
	return err
}

func (dao *dbUserDAO) Exists(id uint) (bool, error) {
	var cnt int
	var err = dao.db.QueryRow(checkUser, id).Scan(&cnt)
	if err != nil {
		return false, err
	}

	return cnt > 0, nil
}
