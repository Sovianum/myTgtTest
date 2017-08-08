package handlers

import (
	"net/http"
)

func ValidateMethod(methodName string, handler handlerType) handlerType {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != methodName {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		handler(w, r)
	}
}

func ValidateNonEmptyBody(errMsg string, handler handlerType) handlerType {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Body == nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(errMsg))
			return
		}

		handler(w, r)
	}
}
