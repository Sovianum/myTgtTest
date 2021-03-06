package handlers

import (
	"github.com/Sovianum/myTgtTest/handlers/mocks"
	"github.com/Sovianum/myTgtTest/model"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	urlSample = "/urlSample"
)

type headerPair struct {
	key   string
	value string
}

func TestEnv_GetRegisterHandler_ContentType(t *testing.T) {
	var env = &Env{Logger: defaultLogger()}
	var rr, err = getRecorder(urlSample, http.MethodPost, env.GetRegisterHandler(), nil)
	if err != nil {
		t.Fatal(err)
	}

	if status := rr.Code; status != http.StatusUnsupportedMediaType {
		t.Errorf("Expected %v, got %v", http.StatusUnsupportedMediaType, status)
	}
}

func TestEnv_GetRegisterHandler_EmptyBody(t *testing.T) {
	var env = &Env{Logger: defaultLogger()}
	var rr, err = getRecorder(
		urlSample,
		http.MethodPost,
		env.GetRegisterHandler(),
		nil,
		headerPair{"Content-Type", "application/json"},
	)
	if err != nil {
		t.Fatal(err)
	}

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Responce code for badly formatted expected %v, got %v", http.StatusBadRequest, status)
	}

	if msg := string(rr.Body.Bytes()); msg != emptyBodyMsg {
		t.Errorf("Message expected \"%v\" \n got \"%v\"", emptyBodyMsg, msg)
	}

}

func TestEnv_GetRegisterHandler_JSONUnparsable(t *testing.T) {
	var env = &Env{Logger: defaultLogger()}
	var rr, err = getRecorder(
		urlSample,
		http.MethodPost,
		env.GetRegisterHandler(),
		strings.NewReader("{it is badly formatted json}"),
		headerPair{"Content-Type", "application/json"},
	)
	if err != nil {
		t.Fatal(err)
	}

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Responce code for badly formatted expected %v, got %v", http.StatusBadRequest, status)
	}
}

func TestEnv_GetRegisterHandler_IncompleteData(t *testing.T) {
	var testData = []struct {
		inputMsg string
		respMsg  string
	}{
		{"{\"age\":100, \"sex\":\"M\"}", model.RegistrationRequiredId},
		{"{\"id\":100, \"sex\":\"M\"}", model.RegistrationRequiredAge},
		{"{\"age\":100, \"id\":100}", model.RegistrationRequiredSex},
	}

	var env = &Env{Logger: defaultLogger()}
	for _, item := range testData {
		var rec, err = getRecorder(
			urlSample,
			http.MethodPost,
			env.GetRegisterHandler(),
			strings.NewReader(item.inputMsg),
			headerPair{"Content-Type", "application/json"},
		)
		if err != nil {
			t.Fatal(err)
		}

		if status := rec.Code; status != http.StatusBadRequest {
			t.Errorf("Id field required. Expected status code %v, got %v on request %v", http.StatusBadRequest, status, item.inputMsg)
		} else if msg := string(rec.Body.Bytes()); msg != item.respMsg {
			t.Errorf("Wrong response expected \n \"%v\" \n, got \n \"%v\" on request \"%v\"", item.respMsg, msg, item.inputMsg)
		}
	}
}

func TestEnv_GetRegisterHandler_IncorrectData(t *testing.T) {
	var testData = []struct {
		inputMsg string
	}{
		{"{\"id\":-100, \"age\":100, \"sex\":\"M\"}"},
		{"{\"id\":\"asdfasdf\", \"age\":100, \"sex\":\"M\"}"},
		{"{\"id\":100, \"age\":-100, \"sex\":\"F\"}"},
		{"{\"id\":100, \"age\":\"as;lkdf\", \"sex\":\"F\"}"},
		{"{\"id\":100, \"age\":100, \"sex\":\"Some\"}"},
		{"{\"id\":100, \"age\":100, \"sex\":90}"},
	}

	var env = &Env{Logger: defaultLogger()}
	for _, item := range testData {
		var rec, err = getRecorder(
			urlSample,
			http.MethodPost,
			env.GetRegisterHandler(),
			strings.NewReader(item.inputMsg),
			headerPair{"Content-Type", "application/json"},
		)
		if err != nil {
			t.Fatal(err)
		}

		if status := rec.Code; status != http.StatusBadRequest {
			t.Errorf("Expected status code %v, got %v on request %v", http.StatusBadRequest, status, item.inputMsg)
		}
	}
}

func TestEnv_GetRegisterHandler_Uniqueness(t *testing.T) {
	var inputMsg = "{\"id\":1, \"age\":1, \"sex\":\"M\"}"
	var successEnv = &Env{UserDAO: new(mocks.NotExistUserDAOMock), Logger: defaultLogger()}
	var failEnv = &Env{UserDAO: new(mocks.ExistUserDAOMock), Logger: defaultLogger()}

	var successRec, successRecErr = getRecorder(
		urlSample,
		http.MethodPost,
		successEnv.GetRegisterHandler(),
		strings.NewReader(inputMsg),
		headerPair{"Content-Type", "application/json"},
	)
	if successRecErr != nil {
		t.Fatal(successRecErr)
	}
	if status := successRec.Code; status != http.StatusOK {
		t.Errorf("Expected status code %v, got %v", http.StatusOK, status)
	}

	var failRec, failRecErr = getRecorder(
		urlSample,
		http.MethodPost,
		failEnv.GetRegisterHandler(),
		strings.NewReader(inputMsg),
		headerPair{"Content-Type", "application/json"},
	)
	if failRecErr != nil {
		t.Fatal(successRecErr)
	}
	if status := failRec.Code; status != http.StatusConflict {
		t.Errorf("Expected status code %v, got %v", http.StatusConflict, status)
	}
}

func TestEnv_GetRegisterHandler_DBError(t *testing.T) {
	var inputMsg = "{\"id\":1, \"age\":1, \"sex\":\"M\"}"
	var env = &Env{UserDAO: new(mocks.FailUserDAOMock), Logger: defaultLogger()}

	var rec, err = getRecorder(
		urlSample,
		http.MethodPost,
		env.GetRegisterHandler(),
		strings.NewReader(inputMsg),
		headerPair{"Content-Type", "application/json"},
	)
	if err != nil {
		t.Fatal(err)
	}
	if status := rec.Code; status != http.StatusInternalServerError {
		t.Errorf("Expected status code %v, got %v", http.StatusInternalServerError, status)
	}
}

func TestEnv_GetStatsAddHandler_EmptyBody(t *testing.T) {
	var env = &Env{Logger: defaultLogger()}
	var rr, err = getRecorder(
		urlSample,
		http.MethodPost,
		env.GetStatsAddHandler(),
		nil,
		headerPair{"Content-Type", "application/json"},
	)
	if err != nil {
		t.Fatal(err)
	}

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Responce code for badly formatted expected %v, got %v", http.StatusBadRequest, status)
	}

	if msg := string(rr.Body.Bytes()); msg != emptyBodyMsg {
		t.Errorf("Message expected \"%v\" \n got \"%v\"", emptyBodyMsg, msg)
	}
}

func TestEnv_GetStatsAddHandler_JSONUnparsable(t *testing.T) {
	var env = &Env{Logger: defaultLogger()}
	var rr, err = getRecorder(
		urlSample,
		http.MethodPost,
		env.GetStatsAddHandler(),
		strings.NewReader("{it is badly formatted json}"),
		headerPair{"Content-Type", "application/json"},
	)
	if err != nil {
		t.Fatal(err)
	}

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Responce code for badly formatted expected %v, got %v", http.StatusBadRequest, status)
	}
}

func TestEnv_GetStatsAddHandler_IncompleteData(t *testing.T) {
	var testData = []struct {
		inputMsg string
		respMsg  string
	}{
		{"{\"action\":\"like\", \"ts\":\"2017-06-30T14:12:34\"}", model.StatsRequiredUser},
		{"{\"user\":100, \"ts\":\"2017-06-30T14:12:34\"}", model.StatsRequiredAction},
		{"{\"user\":100, \"action\":\"like\"}", model.StatsRequiredTs},
	}

	var env = &Env{Logger: defaultLogger()}
	for _, item := range testData {
		var rec, err = getRecorder(
			urlSample,
			http.MethodPost,
			env.GetStatsAddHandler(),
			strings.NewReader(item.inputMsg),
			headerPair{"Content-Type", "application/json"},
		)
		if err != nil {
			t.Fatal(err)
		}

		if status := rec.Code; status != http.StatusBadRequest {
			t.Errorf("Id field required. Expected status code %v, got %v on request %v", http.StatusBadRequest, status, item.inputMsg)
		} else if msg := string(rec.Body.Bytes()); msg != item.respMsg {
			t.Errorf("Wrong response expected \n \"%v\" \n, got \n \"%v\" on request \"%v\"", item.respMsg, msg, item.inputMsg)
		}
	}
}

func TestEnv_GetStatsAddHandler_IncorrectData(t *testing.T) {
	var testData = []struct {
		inputMsg string
	}{
		{"{\"user\":-100, \"action\":\"like\", \"ts\":\"2017-06-30T14:12:34\"}"},
		{"{\"user\":100, \"action\":79, \"ts\":\"2017-06-30T14:12:34\"}"},
		{"{\"user\":100, \"action\":\"like\", \"ts\":\"2017-06-30 14:12:34\"}"},
	}

	var env = &Env{Logger: defaultLogger()}
	for _, item := range testData {
		var rec, err = getRecorder(
			urlSample,
			http.MethodPost,
			env.GetStatsAddHandler(),
			strings.NewReader(item.inputMsg),
			headerPair{"Content-Type", "application/json"},
		)
		if err != nil {
			t.Fatal(err)
		}

		if status := rec.Code; status != http.StatusBadRequest {
			t.Errorf("Expected status code %v, got %v on request %v", http.StatusBadRequest, status, item.inputMsg)
		}
	}
}

func TestEnv_GetStatsAddHandler_UserNotExist(t *testing.T) {
	var inputMsg = "{\"user\":100, \"action\":\"like\", \"ts\":\"2017-06-30T14:12:34\"}"
	var failEnv = &Env{UserDAO: new(mocks.NotExistUserDAOMock), Logger: defaultLogger()}

	var rec, err = getRecorder(
		urlSample,
		http.MethodPost,
		failEnv.GetStatsAddHandler(),
		strings.NewReader(inputMsg),
		headerPair{"Content-Type", "application/json"},
	)
	if err != nil {
		t.Fatal(err)
	}

	if status := rec.Code; status != http.StatusNotFound {
		t.Errorf("Expected status code %v, got %v", http.StatusNotFound, status)
	}
}

func TestEnv_GetStatsAddHandler_UserDAODBError(t *testing.T) {
	var inputMsg = "{\"user\":100, \"action\":\"like\", \"ts\":\"2017-06-30T14:12:34\"}"
	var failEnv = &Env{UserDAO: new(mocks.FailUserDAOMock), Logger: defaultLogger()}

	var rec, err = getRecorder(
		urlSample,
		http.MethodPost,
		failEnv.GetStatsAddHandler(),
		strings.NewReader(inputMsg),
		headerPair{"Content-Type", "application/json"},
	)
	if err != nil {
		t.Fatal(err)
	}

	if status := rec.Code; status != http.StatusInternalServerError {
		t.Errorf("Expected status code %v, got %v", http.StatusInternalServerError, status)
	}
}

func TestEnv_GetStatsAddHandler_StatsDAODBError(t *testing.T) {
	var inputMsg = "{\"user\":100, \"action\":\"like\", \"ts\":\"2017-06-30T14:12:34\"}"
	var failEnv = &Env{
		UserDAO:  new(mocks.ExistUserDAOMock),
		StatsDAO: new(mocks.FailStatsDAOMock),
		Logger:   defaultLogger(),
	}

	var rec, err = getRecorder(
		urlSample,
		http.MethodPost,
		failEnv.GetStatsAddHandler(),
		strings.NewReader(inputMsg),
		headerPair{"Content-Type", "application/json"},
	)
	if err != nil {
		t.Fatal(err)
	}

	if status := rec.Code; status != http.StatusInternalServerError {
		t.Errorf("Expected status code %v, got %v", http.StatusInternalServerError, status)
	}
}

func TestEnv_GetStatsAddHandler_Success(t *testing.T) {
	var inputMsg = "{\"user\":100, \"action\":\"like\", \"ts\":\"2017-06-30T14:12:34\"}"
	var failEnv = &Env{
		UserDAO:  new(mocks.ExistUserDAOMock),
		StatsDAO: new(mocks.SuccessStatsDAOMock),
		Logger:   defaultLogger(),
	}

	var rec, err = getRecorder(
		urlSample,
		http.MethodPost,
		failEnv.GetStatsAddHandler(),
		strings.NewReader(inputMsg),
		headerPair{"Content-Type", "application/json"},
	)
	if err != nil {
		t.Fatal(err)
	}

	if status := rec.Code; status != http.StatusOK {
		t.Errorf("Expected status code %v, got %v", http.StatusOK, status)
	}
}

func TestEnv_GetStatsRequestHandler_IncompleteQueryString(t *testing.T) {
	var urlSlice = []string{
		"/urlSample?action=comments&limit=10",
		"/urlSample?date1=2017-06-20&date2=2017-06-30&limit=10",
		"/urlSample?date1=2017-06-20&date2=2017-06-30&action=comments",
	}

	var env = &Env{Logger: defaultLogger()}
	for _, url := range urlSlice {
		var rec, err = getRecorder(
			url,
			http.MethodGet,
			env.GetStatsRequestHandler(),
			nil,
		)

		if err != nil {
			t.Fatal(err)
		}

		if status := rec.Code; status != http.StatusBadRequest {
			t.Errorf("Expected %v, got %v", http.StatusBadRequest, status)
		}
	}
}

func TestEnv_GetStatsRequestHandler_BadQueryString(t *testing.T) {
	var testData = []struct {
		url            string
		nonPerfectness string
	}{
		{
			url:            "/urlSample?date1=2017-06D20&date2=2017-06-30&action=comments&limit=10",
			nonPerfectness: "Bad formatted date",
		},
		{
			url:            "/urlSample?date1=2017-06-20&date2=2017-06-30&action=coments&limit=10",
			nonPerfectness: "Unknown action",
		},
		{
			url:            "/urlSample?date1=2017-06-20&date2=2017-06-30&action=comments&action=comments&limit=10",
			nonPerfectness: "Many actions",
		},
		{
			url:            "/urlSample?date1=2017-06-20&date2=2017-06-30&action=comments&limit=bad",
			nonPerfectness: "Bad limit value",
		},
		{
			url:            "/urlSample?date1=2017-06-20&date2=2017-06-30&action=comments&limit=10&limit=100",
			nonPerfectness: "Many limits",
		},
	}

	var env = &Env{Logger: defaultLogger()}
	for _, item := range testData {
		var rec, err = getRecorder(
			item.url,
			http.MethodGet,
			env.GetStatsRequestHandler(),
			nil,
		)

		if err != nil {
			t.Fatal(err)
		}

		if status := rec.Code; status != http.StatusBadRequest {
			t.Errorf("Expected %v, got %v (%v)", http.StatusBadRequest, status, item.nonPerfectness)
		}
	}
}

func TestEnv_GetStatsRequestHandler_DBError(t *testing.T) {
	var url = "/urlSample?date1=2017-06-20&date2=2017-06-30&action=comments&limit=10"
	var env = &Env{StatsDAO: new(mocks.FailStatsDAOMock), Logger: defaultLogger()}

	var rec, err = getRecorder(
		url,
		http.MethodGet,
		env.GetStatsRequestHandler(),
		nil,
	)

	if err != nil {
		t.Fatal(err)
	}

	if status := rec.Code; status != http.StatusInternalServerError {
		t.Errorf("Expected %v, got %v", http.StatusInternalServerError, status)
	}
}

func TestEnv_GetStatsRequestHandler_Success(t *testing.T) {
	var testData = []struct {
		url            string
		nonPerfectness string
	}{
		{
			url:            "/urlSample?date1=2017-06-20&date2=2017-06-30&action=comments&limit=10",
			nonPerfectness: "None",
		},
	}

	var env = &Env{StatsDAO: &mocks.SuccessStatsDAOMock{}, Logger: defaultLogger()}

	for _, item := range testData {
		var rec, err = getRecorder(
			item.url,
			http.MethodGet,
			env.GetStatsRequestHandler(),
			nil,
		)

		if err != nil {
			t.Fatal(err)
		}

		if status := rec.Code; status != http.StatusOK {
			t.Errorf("Expected %v, got %v (%v)", http.StatusOK, status, item.nonPerfectness)
		}
	}
}

func getRecorder(url string, method string, handlerFunc HandlerType, body io.Reader, headers ...headerPair) (*httptest.ResponseRecorder, error) {
	var req, err = http.NewRequest(
		method,
		url,
		body,
	)

	for _, hp := range headers {
		req.Header.Set(hp.key, hp.value)
	}

	if err != nil {
		return nil, err
	}

	var handler = http.HandlerFunc(handlerFunc)
	var rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	return rec, nil
}
