package mylog

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/op/go-logging"
	"io"
	"net/http"
	"testing"
)

const (
	method = "METHOD"
	url    = "/URL"
)

func TestLogger_LogBadContentType(t *testing.T) {
	var writer bytes.Buffer
	var logger = getLogger(&writer)
	var req, _ = http.NewRequest(
		method,
		url,
		nil,
	)

	var contentType = "ContentType"
	logger.LogBadContentType(req, contentType)

	var msg = string(writer.Bytes())
	var expected = fmt.Sprintf(
		badContentTypeTemplate,
		req.URL.Path,
		req.Header["Content-Type"],
		contentType,
	)

	if msg[:len(msg)-1] != expected { //msg[:len(msg) - 1] removes last \n symbol
		t.Errorf("Expected \"%v\"; got \"%v\"", expected, msg)
	}
}

func TestLogger_LogRequestStart(t *testing.T) {
	var writer bytes.Buffer
	var logger = getLogger(&writer)
	var req, _ = http.NewRequest(
		method,
		url,
		nil,
	)

	logger.LogRequestStart(req)
	var msg = string(writer.Bytes())
	var expected = fmt.Sprintf(
		requestStartLogTemplate,
		req.URL.Path,
		req.Method,
	)

	if msg[:len(msg)-1] != expected { //msg[:len(msg) - 1] removes last \n symbol
		t.Errorf("Expected \"%v\"; got \"%v\"", expected, msg)
	}
}

func TestLogger_LogEmptyBody(t *testing.T) {
	var writer bytes.Buffer
	var logger = getLogger(&writer)
	var req, _ = http.NewRequest(
		method,
		url,
		nil,
	)

	logger.LogEmptyBody(req)
	var msg = string(writer.Bytes())
	var expected = fmt.Sprintf(
		emptyBodyTemplate,
		req.URL.Path,
		req.Method,
	)

	if msg[:len(msg)-1] != expected { //msg[:len(msg) - 1] removes last \n symbol
		t.Errorf("Expected \"%v\"; got \"%v\"", expected, msg)
	}
}

func TestLogger_LogRequestSuccess(t *testing.T) {
	var writer bytes.Buffer
	var logger = getLogger(&writer)
	var req, _ = http.NewRequest(
		method,
		url,
		nil,
	)

	logger.LogRequestSuccess(req)
	var msg = string(writer.Bytes())
	var expected = fmt.Sprintf(
		requestSuccessLogTemplate,
		req.URL.Path,
		req.Method,
	)

	if msg[:len(msg)-1] != expected { //msg[:len(msg) - 1] removes last \n symbol
		t.Errorf("Expected \"%v\"; got \"%v\"", expected, msg)
	}
}

func TestLogger_LogRequestError(t *testing.T) {
	var writer bytes.Buffer
	var logger = getLogger(&writer)
	var req, _ = http.NewRequest(
		method,
		url,
		nil,
	)

	var errorMsg = "Msg"
	var err = errors.New(errorMsg)

	logger.LogRequestError(req, err)
	var msg = string(writer.Bytes())
	var expected = fmt.Sprintf(
		requestErrorTemplate,
		req.URL.Path,
		err.Error(),
	)

	if msg[:len(msg)-1] != expected { //msg[:len(msg) - 1] removes last \n symbol
		t.Errorf("Expected \"%v\"; got \"%v\"", expected, msg)
	}
}

func TestLogger_LogUserAlreadyExists(t *testing.T) {
	var writer bytes.Buffer
	var logger = getLogger(&writer)
	var req, _ = http.NewRequest(
		method,
		url,
		nil,
	)

	var userid uint = 100

	logger.LogUserAlreadyExists(req, userid)
	var msg = string(writer.Bytes())
	var expected = fmt.Sprintf(
		requestErrorTemplate,
		req.URL.Path,
		fmt.Sprintf(userAlreadyExistsTemplate, userid),
	)

	if msg[:len(msg)-1] != expected { //msg[:len(msg) - 1] removes last \n symbol
		t.Errorf("Expected \"%v\"; got \"%v\"", expected, msg)
	}
}

func TestLogger_LogUserNotExists(t *testing.T) {
	var writer bytes.Buffer
	var logger = getLogger(&writer)
	var req, _ = http.NewRequest(
		method,
		url,
		nil,
	)

	var userid uint = 100

	logger.LogUserNotExists(req, userid)
	var msg = string(writer.Bytes())
	var expected = fmt.Sprintf(
		requestErrorTemplate,
		req.URL.Path,
		fmt.Sprintf(userNotExistsTemplate, userid),
	)

	if msg[:len(msg)-1] != expected { //msg[:len(msg) - 1] removes last \n symbol
		t.Errorf("Expected \"%v\"; got \"%v\"", expected, msg)
	}
}

func getLogger(writer io.Writer) *Logger {
	var format = logging.MustStringFormatter(
		`%{message}`,
	)
	backend := logging.NewLogBackend(writer, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)

	backendLeveled := logging.AddModuleLevel(backendFormatter)
	backendLeveled.SetLevel(logging.INFO, "")

	var logger = logging.MustGetLogger("main")

	logger.SetBackend(backendLeveled)

	return &Logger{*logger}
}
