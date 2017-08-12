package handlers

import (
	"net/http"
)

func LogOnEnter(handler HandlerType, env *Env) HandlerType {
	return func(w http.ResponseWriter, r *http.Request) {
		env.Logger.LogRequestStart(r)
		handler(w, r)
	}
}

// This function is a decorator which takes a request handler and contentType as input
// and returns another handler returning UnsupportedMediaType status code if requested
// with wrong content type header
func ValidateContentType(contentType string, handler HandlerType, env *Env) HandlerType {
	return func(w http.ResponseWriter, r *http.Request) {
		var gotContentType = r.Header.Get("Content-Type")
		if gotContentType != contentType {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			env.Logger.LogBadContentType(r, contentType)
			return
		}

		handler(w, r)
	}
}

// This function is a decorator which takes a request handler error message as input.
// Output handler returns BadRequest status code and error message as response body if request body is empty
func ValidateNonEmptyBody(errMsg string, handler HandlerType, env *Env) HandlerType {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Body == nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(errMsg))
			env.Logger.LogEmptyBody(r)
			return
		}
		handler(w, r)
	}
}
