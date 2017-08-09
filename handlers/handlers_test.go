package handlers

import (
	"github.com/Sovianum/myTgtTest/handlers/mocks"
	"github.com/Sovianum/myTgtTest/model"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"io"
	"strings"
	"github.com/Sovianum/myTgtTest/common"
)

const (
	URL = "/url"
)

func TestEnv_GetRegisterHandler_Method(t *testing.T) {
	log.Println("Started http method testing")

	var getRR, getErr = getRecorder(URL, http.MethodGet, new(Env).GetRegisterHandler(), nil)
	if getErr != nil {
		t.Fatal(getErr)
	}

	if status := getRR.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("Method get not allowed for this request: got %v expected %v", status, http.StatusMethodNotAllowed)
	}

	log.Println("Http method tested successfully")
}

func TestEnv_GetRegisterHandler_EmptyBody(t *testing.T) {
	log.Println("Started empty body testing")

	var rr, err = getRecorder(URL, http.MethodPost, new(Env).GetRegisterHandler(), nil)
	if err != nil {
		t.Fatal(err)
	}

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Responce code for badly formatted expected %v, got %v", http.StatusBadRequest, status)
	}

	if msg := string(rr.Body.Bytes()); msg != emptyBodyMsg {
		t.Errorf("Message expected \"%v\" \n got \"%v\"", emptyBodyMsg, msg)
	}

	log.Println("Empty body tested successfully")

}

func TestEnv_GetRegisterHandler_JSONUnparsable(t *testing.T) {
	log.Println("Started unparsable json testing")

	var rr, err = getRecorder(
		URL,
		http.MethodPost,
		new(Env).GetRegisterHandler(),
		strings.NewReader("{it is badly formatted json}"),
	)
	if err != nil {
		t.Fatal(err)
	}

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Responce code for badly formatted expected %v, got %v", http.StatusBadRequest, status)
	}

	log.Println("JSON parsing tested successfully")
}

func TestEnv_GetRegisterHandler_IncompleteData(t *testing.T) {
	log.Println("Started incomplete json testing")

	var testData = []struct {
		inputMsg string
		respMsg  string
	}{
		{"{\"age\":100, \"sex\":\"M\"}", model.RegistrationRequiredId},
		{"{\"id\":100, \"sex\":\"M\"}", model.RegistrationRequiredAge},
		{"{\"age\":100, \"id\":100}", model.RegistrationRequiredSex},
	}

	for _, item := range testData {
		var rec, err = getRecorder(
			URL,
			http.MethodPost,
			new(Env).GetRegisterHandler(),
			strings.NewReader(item.inputMsg),
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

	log.Println("Incomplete json tested successfully")
}

func TestEnv_GetRegisterHandler_IncorrectData(t *testing.T) {
	log.Println("Started incorrect json testing")

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

	for _, item := range testData {
		var rec, err = getRecorder(
			URL,
			http.MethodPost,
			new(Env).GetRegisterHandler(),
			strings.NewReader(item.inputMsg),
		)
		if err != nil {
			t.Fatal(err)
		}

		if status := rec.Code; status != http.StatusBadRequest {
			t.Errorf("Expected status code %v, got %v on request %v", http.StatusBadRequest, status, item.inputMsg)
		}
	}

	log.Println("Incorrect json tested successfully")
}

func TestEnv_GetRegisterHandler_Uniqueness(t *testing.T) {
	log.Println("Started user uniqueness testing")
	var inputMsg = "{\"id\":1, \"age\":1, \"sex\":\"M\"}"
	var successEnv = &Env{userDAO: new(mocks.NotExistUserDAOMock)}
	var failEnv = &Env{userDAO: new(mocks.ExistUserDAOMock)}

	var successRec, successRecErr = getRecorder(
		URL,
		http.MethodPost,
		successEnv.GetRegisterHandler(),
		strings.NewReader(inputMsg),
	)
	if successRecErr != nil {
		t.Fatal(successRecErr)
	}
	if status := successRec.Code; status != http.StatusOK {
		t.Errorf("Expected status code %v, got %v", http.StatusOK, status)
	}

	var failRec, failRecErr = getRecorder(
		URL,
		http.MethodPost,
		failEnv.GetRegisterHandler(),
		strings.NewReader(inputMsg),
	)
	if failRecErr != nil {
		t.Fatal(successRecErr)
	}
	if status := failRec.Code; status != http.StatusConflict {
		t.Errorf("Expected status code %v, got %v", http.StatusConflict, status)
	}

	log.Println("User uniqueness tested successfully")
}

func TestEnv_GetStatsAddHandler_Method(t *testing.T) {
	log.Println("Started http method testing")

	var getRR, getErr = getRecorder(URL, http.MethodGet, new(Env).GetStatsAddHandler(), nil)
	if getErr != nil {
		t.Fatal(getErr)
	}

	if status := getRR.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("Method get not allowed for this request: got %v expected %v", status, http.StatusMethodNotAllowed)
	}

	log.Println("Http method tested successfully")
}

func TestEnv_GetStatsAddHandler_EmptyBody(t *testing.T) {
	log.Println("Started empty body testing")

	var rr, err = getRecorder(URL, http.MethodPost, new(Env).GetStatsAddHandler(), nil)
	if err != nil {
		t.Fatal(err)
	}

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Responce code for badly formatted expected %v, got %v", http.StatusBadRequest, status)
	}

	if msg := string(rr.Body.Bytes()); msg != emptyBodyMsg {
		t.Errorf("Message expected \"%v\" \n got \"%v\"", emptyBodyMsg, msg)
	}

	log.Println("Empty body tested successfully")

}

func TestEnv_GetStatsAddHandler_JSONUnparsable(t *testing.T) {
	log.Println("Started unparsable json testing")

	var rr, err = getRecorder(
		URL,
		http.MethodPost,
		new(Env).GetStatsAddHandler(),
		strings.NewReader("{it is badly formatted json}"),
	)
	if err != nil {
		t.Fatal(err)
	}

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Responce code for badly formatted expected %v, got %v", http.StatusBadRequest, status)
	}

	log.Println("JSON parsing tested successfully")
}

func TestEnv_GetStatsAddHandler_IncompleteData(t *testing.T) {
	log.Println("Started incomplete json testing")

	var testData = []struct {
		inputMsg string
		respMsg  string
	}{
		{"{\"action\":\"like\", \"ts\":\"2017-06-30T14:12:34\"}", model.StatsRequiredUser},
		{"{\"user\":100, \"ts\":\"2017-06-30T14:12:34\"}", model.StatsRequiredAction},
		{"{\"user\":100, \"action\":\"like\"}", model.StatsRequiredTs},
	}

	for _, item := range testData {
		var rec, err = getRecorder(
			URL,
			http.MethodPost,
			new(Env).GetStatsAddHandler(),
			strings.NewReader(item.inputMsg),
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

	log.Println("Incomplete json tested successfully")
}

func TestEnv_GetStatsAddHandler_IncorrectData(t *testing.T) {
	log.Println("Started incorrect json testing")

	var testData = []struct {
		inputMsg string
	}{
		{"{\"user\":-100, \"action\":\"like\", \"ts\":\"2017-06-30T14:12:34\"}"},
		{"{\"user\":100, \"action\":79, \"ts\":\"2017-06-30T14:12:34\"}"},
		{"{\"user\":100, \"action\":\"like\", \"ts\":\"2017-06-30 14:12:34\"}"},
	}

	for _, item := range testData {
		var rec, err = getRecorder(
			URL,
			http.MethodPost,
			new(Env).GetStatsAddHandler(),
			strings.NewReader(item.inputMsg),
		)
		if err != nil {
			t.Fatal(err)
		}

		if status := rec.Code; status != http.StatusBadRequest {
			t.Errorf("Expected status code %v, got %v on request %v", http.StatusBadRequest, status, item.inputMsg)
		}
	}

	log.Println("Incorrect json tested successfully")
}

func TestEnv_GetStatsAddHandler_UserNotExist(t *testing.T) {
	log.Println("Started user does not exist testing")
	var inputMsg = "{\"user\":100, \"action\":\"like\", \"ts\":\"2017-06-30T14:12:34\"}"
	var failEnv = &Env{userDAO: new(mocks.NotExistUserDAOMock)}

	var rec, err = getRecorder(
		URL,
		http.MethodPost,
		failEnv.GetStatsAddHandler(),
		strings.NewReader(inputMsg),
	)
	if err != nil {
		t.Fatal(err)
	}

	if status := rec.Code; status != http.StatusNotFound {
		t.Errorf("Expected status code %v, got %v", http.StatusNotFound, status)
	}

	log.Println("User does not exist tested successfully")
}

func TestEnv_GetStatsAddHandler_Success(t *testing.T) {
	log.Println("Started success testing")
	var inputMsg = "{\"user\":100, \"action\":\"like\", \"ts\":\"2017-06-30T14:12:34\"}"
	var failEnv = &Env{userDAO: new(mocks.ExistUserDAOMock), statsDAO:new(mocks.SuccessStatsDaoMock)}

	var rec, err = getRecorder(
		URL,
		http.MethodPost,
		failEnv.GetStatsAddHandler(),
		strings.NewReader(inputMsg),
	)
	if err != nil {
		t.Fatal(err)
	}

	if status := rec.Code; status != http.StatusOK {
		t.Errorf("Expected status code %v, got %v", http.StatusOK, status)
	}

	log.Println("Success tested successfully")
}

func TestEnv_GetStatsRequestHandler_Method(t *testing.T) {
	log.Println("Started http method testing")

	var rec, err = getRecorder(URL, http.MethodPost, new(Env).GetStatsRequestHandler(), nil)
	if err != nil {
		t.Fatal(err)
	}

	if status := rec.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("Method get not allowed for this request: got %v expected %v", status, http.StatusMethodNotAllowed)
	}

	log.Println("Http method tested successfully")
}

//func TestEnv_GetStatsRequestHandler_IncompleteQueryString(t *testing.T) {
//	log.Println("Started incomplete query string testing")
//
//	var rec, err = getRecorder("/url?date1=2017-06-20", http.MethodPost, new(Env).GetStatsRequestHandler(), nil)
//
//	log.Println("Incomplete query string tested successfully")
//}

func getRecorder(url string, method string, handlerFunc common.HandlerType, body io.Reader) (*httptest.ResponseRecorder, error) {
	var req, err = http.NewRequest(
		method,
		url,
		body,
	)

	if err != nil {
		return nil, err
	}

	var handler = http.HandlerFunc(handlerFunc)
	var rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	return rec, nil
}