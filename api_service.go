package apid

import "net/http"

type APIService interface {
	Listen() error
	Handle(path string, handler http.Handler) Route
	HandleFunc(path string, handlerFunc http.HandlerFunc) Route
	Vars(r *http.Request) map[string]string
	Router() Router
}

type Router interface {
	Handle(path string, handler http.Handler) Route
	HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) Route
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

type Route interface {
	Methods(methods ...string) Route
}
