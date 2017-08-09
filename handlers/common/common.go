package common

import "net/http"

type HandlerType func(http.ResponseWriter, *http.Request)
