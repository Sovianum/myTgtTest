package handlers

import (
	"net/http"
	"strings"
	"testing"
)

func TestValidateMethod_Success(t *testing.T) {
	var innerFunc = func(w http.ResponseWriter, r *http.Request) {}
	var wrappedFunc = ValidateMethod(http.MethodGet, innerFunc)

	var rec, _ = getRecorder(urlSample, http.MethodGet, wrappedFunc, nil)

	if status := rec.Code; status != http.StatusOK {
		t.Errorf("Expected %v, got %v", http.StatusOK, status)
	}
}

func TestValidateMethod_Fail(t *testing.T) {
	var innerFunc = func(w http.ResponseWriter, r *http.Request) {}
	var wrappedFunc = ValidateMethod(http.MethodGet, innerFunc)

	var rec, _ = getRecorder(urlSample, http.MethodPost, wrappedFunc, nil)

	if status := rec.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("Expected %v, got %v", http.StatusMethodNotAllowed, status)
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
