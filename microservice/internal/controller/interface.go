package controller

import "net/http"

type HTTPHandler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}
