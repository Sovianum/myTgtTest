package handlers

import (
	"net/http"
	"strings"
	"testing"
)

func TestValidateContentType_Fail(t *testing.T) {
	var innerFunc = func(w http.ResponseWriter, r *http.Request) {}
	var wrappedFunc = ValidateContentType("application/json", innerFunc)

	var rec, _ = getRecorder(urlSample, http.MethodPost, wrappedFunc, nil)

	if rec.Code != http.StatusUnsupportedMediaType {
		t.Errorf("Expected %v, got %v", http.StatusUnsupportedMediaType, rec.Code)
	}
}

func TestValidateContentType_Success(t *testing.T) {
	var innerFunc = func(w http.ResponseWriter, r *http.Request) {}
	var wrappedFunc = ValidateContentType("application/json", innerFunc)

	var rec, _ = getRecorder(
		urlSample,
		http.MethodPost, wrappedFunc,
		nil,
		headerPair{"Content-Type", "application/json"},
	)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected %v, got %v", http.StatusOK, rec.Code)
	}
}

func TestValidateNonEmptyBody_Success(t *testing.T) {
	var innerFunc = func(w http.ResponseWriter, r *http.Request) {}
	var wrappedFunc = ValidateNonEmptyBody(http.MethodGet, innerFunc)

	var rec, _ = getRecorder(urlSample, http.MethodGet, wrappedFunc, strings.NewReader(";lkj;l"))

	if status := rec.Code; status != http.StatusOK {
		t.Errorf("Expected %v, got %v", http.StatusOK, status)
	}
}

func TestValidateNonEmptyBody_Fail(t *testing.T) {
	var innerFunc = func(w http.ResponseWriter, r *http.Request) {}
	var wrappedFunc = ValidateNonEmptyBody(http.MethodGet, innerFunc)

	var rec, _ = getRecorder(urlSample, http.MethodGet, wrappedFunc, nil)

	if status := rec.Code; status != http.StatusBadRequest {
		t.Errorf("Expected %v, got %v", http.StatusBadRequest, status)
	}
}
