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
var requests *expvar.Map = expvar.NewMap("requests")

func CreateService() apid.APIService {
	config = apid.Config()
	log = apid.Log().ForModule("api")

	config.SetDefault(configAPIPort, 9000)

	r := mux.NewRouter()
	rw := &router{r}
	return &service{rw}
}

type service struct {
	*router
}

func (s *service) Listen() error {
	port := config.GetString(configAPIPort)
	log.Infof("opening api port %s", port)

	s.InitExpVar()

	apid.Events().Emit(apid.SystemEventsSelector, apid.APIListeningEvent) // todo: run after successful listen?
	return http.ListenAndServe(":"+port, s.r)
}

func (s *service) InitExpVar() {
	if config.IsSet(configExpVarPath) {
		log.Infof("expvar available on path: %s", config.Get(configExpVarPath))
		s.HandleFunc(config.GetString(configExpVarPath), expvarHandler)
	}
}

// for testing
func (s *service) Router() apid.Router {
	s.InitExpVar()
	return s
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
	log.Infof("Handle %s: %v", path, handler)
	return &route{r.r.Handle(path, handler)}
}

func (r *router) HandleFunc(path string, handlerFunc http.HandlerFunc) apid.Route {
	log.Infof("Handle %s: %v", path, handlerFunc)
	return &route{r.r.HandleFunc(path, handlerFunc)}
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	requests.Add(req.URL.Path, 1)
	r.r.ServeHTTP(w, req)
}

type route struct {
	r *mux.Route
}

func (r *route) Methods(methods ...string) apid.Route {
	return &route{r.r.Methods(methods...)}
}
