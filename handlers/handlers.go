package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Sovianum/myTgtTest/dao"
	"github.com/Sovianum/myTgtTest/model"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	userNotFound      = "\"User not found\""
	userAlreadyExists = "\"User already exists\""
	emptyBodyMsg      = "\"Empty body not allowed\""
	requiredActionMsg = "\"Required \"action\" query parameter\""
	requiredLimitMsg  = "\"Required \"limit\" query parameter\""
	requiredDateMsg   = "\"Required at least on date value\""

	badLimitValueMsg  = "\"Used incorrect limit value\""
	badActionValueMsg = "\"Used incorrect action value\""

	manyLimitValuesMsg  = "\"Can not use multiple limit values\""
	manyActionValuesMsg = "\"Can not use multiple action values\""

	datePrefix      = "date"
	actionParameter = "action"
	limitParameter  = "limit"
)

type HandlerType func(http.ResponseWriter, *http.Request)

type Env struct {
	userDAO  dao.UserDAO
	statsDAO dao.StatsDAO
}

func NewDBEnv(db *sql.DB) Env {
	return Env{
		userDAO:  dao.NewDBUserDAO(db),
		statsDAO: dao.NewDBStatsDAO(db),
	}
}

func (env *Env) GetRegisterHandler() HandlerType {
	var innerFunc = func(w http.ResponseWriter, r *http.Request) {
		var registration = model.Registration{}
		var err = registration.ReadJsonIn(r.Body)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		var exists, existsError = env.userDAO.Exists(registration.Id)
		if existsError != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(existsError.Error()))
		}

		if exists {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(userAlreadyExists))
			return
		}

		var saveError = env.userDAO.Save(registration)
		if saveError != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(saveError.Error()))
			return
		}
	}

	return ValidateContentType(
		"application/json",
		ValidateNonEmptyBody(
			emptyBodyMsg,
			innerFunc,
		),
	)
}

func (env *Env) GetStatsAddHandler() HandlerType {
	var innerFunc = func(w http.ResponseWriter, r *http.Request) {
		var stats = model.Stats{}
		var err = stats.ReadJsonIn(r.Body)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		var exists, existsErr = env.userDAO.Exists(stats.User)
		if existsErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(existsErr.Error()))
		}

		if !exists {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(userNotFound))
			return
		}

		var saveError = env.statsDAO.Save(stats)
		if saveError != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(saveError.Error()))
			return
		}
	}

	return ValidateContentType(
		"application/json",
		ValidateNonEmptyBody(
			emptyBodyMsg,
			innerFunc,
		),
	)
}

func (env *Env) GetStatsRequestHandler() HandlerType {
	return func(w http.ResponseWriter, r *http.Request) {
		var query = r.URL.Query()
		var checkErr = checkFieldsExistence(
			query,
			[]string{actionParameter, limitParameter},
			[]string{requiredActionMsg, requiredLimitMsg},
		)

		if checkErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(checkErr.Error()))
			return
		}

		var dateSlice, dateParseErr = getParsedDateSlice(getDateSlice(query))
		if dateParseErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(dateParseErr.Error()))
			return
		}

		var action, limit, parseErr = getActionAndLimit(query)
		if parseErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(parseErr.Error()))
			return
		}

		var statsSlice, err = env.statsDAO.GetStatsSlice(dateSlice, action, limit)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		var msg, _ = json.Marshal(statsSlice)
		w.Write(msg)
	}
}

// Function returns error if some of names in nameSlice not present in query.
// Resulting error message consists of corresponding messages from errMsgSlice joined with ";\n"
func checkFieldsExistence(query url.Values, nameSlice []string, errMsgSlice []string) error {
	var msgSlice = make([]string, 0)

	for i, name := range nameSlice {
		var _, ok = query[name]

		if !ok {
			msgSlice = append(msgSlice, errMsgSlice[i])
		}
	}

	if len(msgSlice) > 0 {
		return errors.New(strings.Join(msgSlice, ";\n"))
	}

	return nil
}

func getActionAndLimit(query url.Values) (string, int, error) {
	var msgSlice = make([]string, 0)

	var limitStrSlice, okLimit = query[limitParameter]
	if !okLimit {
		msgSlice = append(msgSlice, requiredLimitMsg)
	}
	if len(limitStrSlice) > 1 {
		msgSlice = append(msgSlice, manyLimitValuesMsg)
	}

	var limit, limitErr = strconv.Atoi(limitStrSlice[0])
	if limitErr != nil {
		msgSlice = append(msgSlice, badLimitValueMsg, limitErr.Error())
	}

	var actions, okAction = query[actionParameter]
	if !okAction {
		msgSlice = append(msgSlice, requiredActionMsg)
	}
	if len(actions) > 1 {
		msgSlice = append(msgSlice, manyActionValuesMsg)
	}
	if !model.IsValidAction(actions[0]) {
		msgSlice = append(msgSlice, badActionValueMsg)
	}

	if len(msgSlice) > 0 {
		return "", 0, errors.New(strings.Join(msgSlice, ";\n"))
	}

	return actions[0], limit, nil
}

// Function parses output of gatDateSlice function. Empty slice is considered to be an error
func getParsedDateSlice(strDateSlice []string) ([]time.Time, error) {
	if len(strDateSlice) == 0 {
		return []time.Time{}, errors.New(requiredDateMsg)
	}

	var result = make([]time.Time, 0)
	var date time.Time
	var err error

	for _, s := range strDateSlice {
		date, err = time.Parse("2006-01-02", s)
		if err != nil {
			break
		}

		result = append(result, date)
	}

	return result, err
}

// Function returns slice of string values in query, trying sequentially names of type
// date%d, starting from date1. If name is not found, search ends. Empty slice is
// not considered to be an error
func getDateSlice(query url.Values) []string {
	var result = make([]string, 0)
	var dateNameGen = getNameGen(datePrefix, 1)

	for {
		var date, ok = query[dateNameGen()]
		if !ok {
			break
		}

		result = append(result, date...)
	}

	return result
}

// Function returns another function, which generates sequential names
// For example: prefix=date, startNum=1 => date1, date2, date3...
func getNameGen(prefix string, startNum int) func() string {
	var cnt = startNum - 1
	return func() string {
		cnt++
		return fmt.Sprintf("%s%d", prefix, cnt)
	}
}
