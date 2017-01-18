package api

import (
	"github.com/30x/apid"
	"github.com/gorilla/mux"
	"net/http"
	"expvar"
	"fmt"
	"github.com/30x/goscaffold"
)

// todo: handle TLS
// todo: handle other basic router config, errors, etc.

const (
	configAPIPort = "api_port"
	configExpVarPath = "api_expvar_path"
	configReadyPath = "api_ready"
	configHealthyPath = "api_healthy"
)

var log apid.LogService
var config apid.ConfigService
var requests *expvar.Map = expvar.NewMap("requests")

func CreateService() apid.APIService {
	config = apid.Config()
	log = apid.Log().ForModule("api")

	config.SetDefault(configAPIPort, 9000)
	config.SetDefault(configReadyPath, "/ready")
	config.SetDefault(configHealthyPath, "/healthy")

	r := mux.NewRouter()
	rw := &router{r}
	scaffold := goscaffold.CreateHTTPScaffold()
	return &service{rw, scaffold}
}

type service struct {
	*router
	scaffold *goscaffold.HTTPScaffold
}

func (s *service) Listen() error {
	port := config.GetInt(configAPIPort)
	log.Infof("opening api port %d", port)

	s.InitExpVar()

	s.scaffold.SetInsecurePort(port)

	// Direct the scaffold to catch common signals and trigger a graceful shutdown.
	s.scaffold.CatchSignals()

	// Set an URL that may be used by a load balancer to test if the server is ready to handle requests
	if config.GetString(configReadyPath) != "" {
		s.scaffold.SetReadyPath(config.GetString(configReadyPath))
	}

	// Set an URL that may be used by infrastructure to test
	// if the server is working or if it needs to be restarted or replaced
	if config.GetString(configHealthyPath) != "" {
		s.scaffold.SetReadyPath(config.GetString(configHealthyPath))
	}

	err := s.scaffold.StartListen(s.r)
	if err != nil {
		return err
	}

	apid.Events().Emit(apid.SystemEventsSelector, apid.APIListeningEvent)

	return s.scaffold.WaitForShutdown()
}

func (s *service) Close() {
	s.scaffold.Shutdown(nil)
	s.scaffold = nil
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
	log.Infof("Handling %s", req.URL.Path)
	r.r.ServeHTTP(w, req)
}

type route struct {
	r *mux.Route
}

func (r *route) Methods(methods ...string) apid.Route {
	return &route{r.r.Methods(methods...)}
}
