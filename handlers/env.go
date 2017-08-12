package handlers

import (
	"database/sql"
	"github.com/Sovianum/myTgtTest/dao"
	"github.com/Sovianum/myTgtTest/mylog"
)

type Env struct {
	UserDAO  dao.UserDAO
	StatsDAO dao.StatsDAO
	Logger   *mylog.Logger
}

func NewDBEnv(db *sql.DB, logger *mylog.Logger) Env {
	return Env{
		UserDAO:  dao.NewDBUserDAO(db),
		StatsDAO: dao.NewDBStatsDAO(db),
		Logger:   logger,
	}
}
