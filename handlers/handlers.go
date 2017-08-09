package handlers

import (
	"github.com/Sovianum/myTgtTest/dao"
	"github.com/Sovianum/myTgtTest/model"
	"net/http"
	"github.com/Sovianum/myTgtTest/decorators"
	"github.com/Sovianum/myTgtTest/common"
)

const (
	emptyBodyMsg = "\"Empty body not allowed\""
)

type Env struct {
	userDAO  dao.UserDAO
	statsDAO dao.StatsDAO
}

func (env *Env) GetRegisterHandler() common.HandlerType {
	var innerFunc = func(w http.ResponseWriter, r *http.Request) {
		var registration = model.Registration{}
		var err = registration.UnmarshalJSON(r.Body)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		var saveError = env.userDAO.Save(registration)
		if saveError != nil {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(saveError.Error()))
			return
		}
	}

	return decorators.ValidateMethod(
		http.MethodPost,
		decorators.ValidateNonEmptyBody(
			emptyBodyMsg,
			innerFunc,
		),
	)
}

func (env *Env) GetStatsAddHandler() common.HandlerType {
	var innerFunc = func(w http.ResponseWriter, r *http.Request) {
		var stats = model.Stats{}
		var err = stats.UnmarshalJSON(r.Body)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		if !env.userDAO.Exists(stats.User) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		var saveError = env.statsDAO.Save(stats)
		if saveError != nil {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(saveError.Error()))
			return
		}
	}

	return decorators.ValidateMethod(
		http.MethodPost,
		decorators.ValidateNonEmptyBody(
			emptyBodyMsg,
			innerFunc,
		),
	)
}

func (env *Env) GetStatsRequestHandler() common.HandlerType {
	var innerFunc = func(w http.ResponseWriter, r *http.Request) {}

	return decorators.ValidateMethod(
		http.MethodGet,
		innerFunc,
	)
}
