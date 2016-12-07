package apid

import "net/http"

type APIService interface {
	Listen() error
	Handle(path string, handler http.Handler) Route
	HandleFunc(path string, handlerFunc http.HandlerFunc) Route
	Vars(r *http.Request) map[string]string

	// for testing
	Router() Router
}

type Route interface {
	Methods(methods ...string) Route
}

// for testing
type Router interface {
	Handle(path string, handler http.Handler) Route
	HandleFunc(path string, handlerFunc http.HandlerFunc) Route
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}
