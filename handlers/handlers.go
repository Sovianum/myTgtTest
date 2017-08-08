package handlers

import (
	"github.com/Sovianum/myTgtTest/dao"
	"github.com/Sovianum/myTgtTest/model"
	"net/http"
)

type handlerType func(http.ResponseWriter, *http.Request)

type Env struct {
	userDAO  dao.UserDAO
	statsDAO dao.StatsDAO
}

func (env *Env) GetRegisterHandler() handlerType {
	var innerFunc = func(w http.ResponseWriter, r *http.Request) {
		var registration = *new(model.Registration)
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

	return ValidateMethod(
		http.MethodPost,
		ValidateNonEmptyBody(
			emptyBodyMsg,
			innerFunc,
		),
	)
}

func (env *Env) GetStatsAddHandler() handlerType {
	var innerFunc = func(w http.ResponseWriter, r *http.Request) {

	}

	return ValidateMethod(
		http.MethodPost,
		ValidateNonEmptyBody(
			emptyBodyMsg,
			innerFunc,
		),
	)
}
