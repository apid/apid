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

	router := mux.NewRouter()
	return &service{router}
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

func (s *service) Handle(path string, handler http.Handler) {
	log.Infof("handle %s: %s", path, handler)
	s.router.Handle(path, handler)
}

func (s *service) HandleFunc(path string, handlerFunc http.HandlerFunc) {
	log.Infof("handle %s: %s", path, handlerFunc)
	s.router.HandleFunc(path, handlerFunc)
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
