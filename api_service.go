package apid

import "net/http"

type APIService interface {
	Listen() error
	Handle(path string, handler http.Handler)
	HandleFunc(path string, handlerFunc http.HandlerFunc)
}
