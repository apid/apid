package api

import (
	"github.com/30x/apid"
	"github.com/gorilla/mux"
	"net/http"
	"expvar"
	"fmt"
)

// todo: handle TLS
// todo: handle other basic router config, errors, etc.

const (
	configAPIPort = "api_port"
	configExpVarPath = "api_expvar_path"
)

var log apid.LogService
var config apid.ConfigService

func CreateService() apid.APIService {
	config = apid.Config()
	log = apid.Log()

	config.SetDefault(configAPIPort, 9000)

	r := mux.NewRouter()
	return &service{r}
}

type service struct {
	router *mux.Router
}

func (s *service) Listen() error {
	port := config.GetString(configAPIPort)
	log.Infof("opening api port %s", port)

	if config.IsSet(configExpVarPath) {
		log.Info("expvar available on path: %s", config.Get(configExpVarPath))
		s.HandleFunc(config.GetString(configExpVarPath), expvarHandler)
	}

	apid.Events().Emit(apid.SystemEventsSelector, apid.APIListeningEvent) // todo: run after successful listen?
	return http.ListenAndServe(":"+port, s.router)
}

func (s *service) Handle(path string, handler http.Handler) apid.Route {
	log.Infof("handle %s: %s", path, handler)
	return s.Router().Handle(path, handler)
}

func (s *service) HandleFunc(path string, handlerFunc http.HandlerFunc) apid.Route {
	log.Infof("handle %s: %s", path, handlerFunc)
	return s.Router().HandleFunc(path, handlerFunc)
}

func (s *service) Router() apid.Router {
	return &router{s.router}
}

func (s *service) Vars(r *http.Request) map[string]string {
	return mux.Vars(r)
}

func expvarHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprint(w, "{\n")
	first := true
	expvar.Do(func(kv expvar.KeyValue) {
		if !first {
			fmt.Fprint(w, ",\n")
		}
		first = false
		fmt.Fprintf(w, "%q: %s", kv.Key, kv.Value)
	})
	fmt.Fprint(w, "\n}\n")
}

type router struct {
	r *mux.Router
}

func (r *router) Handle(path string, handler http.Handler) apid.Route {
	return &route{r.r.Handle(path, handler)}
}

func (r *router) HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) apid.Route {
	return &route{r.r.HandleFunc(path, f)}
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.r.ServeHTTP(w, req)
}

type route struct {
	r *mux.Route
}

func (r *route) Methods(methods ...string) apid.Route {
	return &route{r.r.Methods(methods...)}
}
