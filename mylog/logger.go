package mylog

import (
	"fmt"
	golog "github.com/op/go-logging"
	"net/http"
)

const (
	requestStartLogTemplate   = `Started handling request to url %v with method %v`
	requestSuccessLogTemplate = `Request to url %v with method %v handled successfully`
	emptyBodyTemplate         = `Request to url %v with method %v has empty body`
	badContentTypeTemplate    = `Request to url %v has Content-Type %v (expected %v)`
	requestErrorTemplate      = `Failed on URL %v with error \"%v\"`
	userAlreadyExistsTemplate = `User with id = %d already exists`
	userNotExistsTemplate     = `User with id = %d not exists`
)

type Logger struct {
	golog.Logger
}

func (logger *Logger) LogRequestStart(r *http.Request) {
	logger.Infof(requestStartLogTemplate, r.URL.Path, r.Method)
}

func (logger *Logger) LogBadContentType(r *http.Request, expectedContentType string) {
	logger.Errorf(badContentTypeTemplate, r.URL.Path, r.Header["Content-Type"], expectedContentType)
}

func (logger *Logger) LogEmptyBody(r *http.Request) {
	logger.Errorf(emptyBodyTemplate, r.URL.Path, r.Method)
}

func (logger *Logger) LogRequestSuccess(r *http.Request) {
	logger.Infof(requestSuccessLogTemplate, r.URL.Path, r.Method)
}

func (logger *Logger) LogRequestError(r *http.Request, err error) {
	logger.Errorf(requestErrorTemplate, r.URL.Path, err.Error())
}

func (logger *Logger) LogUserAlreadyExists(r *http.Request, userId uint) {
	logger.Errorf(
		requestErrorTemplate,
		r.URL.Path,
		fmt.Sprintf(userAlreadyExistsTemplate, userId),
	)
}

func (logger *Logger) LogUserNotExists(r *http.Request, userId uint) {
	logger.Errorf(
		requestErrorTemplate,
		r.URL.Path,
		fmt.Sprintf(userNotExistsTemplate, userId),
	)
}
