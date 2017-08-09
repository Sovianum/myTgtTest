package handlers

import (
	"net/http"
)

// This function is a decorator which takes a request handler and httpMethod as input
// and returns another handler returning MethodNotAllowed status code if requested with another method
func ValidateMethod(methodName string, handler HandlerType) HandlerType {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != methodName {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		handler(w, r)
	}
}

// This function is a decorator which takes a request handler error message as input.
// Output handler returns BadRequest status code and error message as response body if request body is empty
func ValidateNonEmptyBody(errMsg string, handler HandlerType) HandlerType {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Body == nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(errMsg))
			return
		}
		handler(w, r)
	}
}
