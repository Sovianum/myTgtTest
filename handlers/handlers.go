package handlers

import (
	"github.com/Sovianum/myTgtTest/common"
	"github.com/Sovianum/myTgtTest/dao"
	"github.com/Sovianum/myTgtTest/decorators"
	"github.com/Sovianum/myTgtTest/model"
	"net/http"
	"fmt"
	"net/url"
	"time"
	"errors"
	"strings"
	"strconv"
	"encoding/json"
)

const (
	emptyBodyMsg = "\"Empty body not allowed\""
	requiredActionMsg = "\"Required \"action\" query parameter\""
	requiredLimitMsg = "\"Required \"limit\" query parameter\""
	requiredDateMsg = "\"Required at least on date value\""

	badLimitValueMsg = "\"Used incorrect limit value\""
	badActionValueMsg = "\"Used incorrect action value\""

	manyLimitValuesMsg = "\"Can not use multiple limit values\""
	manyActionValuesMsg = "\"Can not use multiple action values\""

	datePrefix = "date"
	actionParameter = "action"
	limitParameter = "limit"
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
	var innerFunc = func(w http.ResponseWriter, r *http.Request) {
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

		var statsSlice, err = env.statsDAO.Get(dateSlice, action, limit)
		if parseErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		var msg, _ = json.Marshal(statsSlice)
		w.Write(msg)
	}

	return decorators.ValidateMethod(
		http.MethodGet,
		innerFunc,
	)
}

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

func getNameGen(prefix string, startNum int) func() string {
	var cnt = startNum - 1
	return func() string {
		cnt++
		return fmt.Sprintf("%s%d", prefix, cnt)
	}
}
