package apid

const (
	SystemEventsSelector EventSelector = "system event"
)

var (
	APIDInitializedEvent        = systemEvent{"apid initialized"}
	PluginsInitializedEvent     = systemEvent{"plugins initialized"}
	PostPluginsInitializedEvent = systemEvent{"post plugins initialized"}
	APIListeningEvent           = systemEvent{"api listening"}

	pluginInitFuncs     []PluginInitFunc
	postpluginInitFuncs []PostPluginInitFunc
	services            Services
)

type Services interface {
	API() APIService
	Config() ConfigService
	Data() DataService
	Events() EventsService
	Log() LogService
}

type PluginInitFunc func(Services) error
type PostPluginInitFunc func(Services) error

// passed Services can be a factory - makes copies and maintains returned references
// eg. apid.Initialize(factory.DefaultServicesFactory())

func Initialize(s Services) {
	ss := &servicesSet{}
	services = ss
	// order is important
	ss.config = s.Config()
	ss.log = s.Log()
	ss.events = s.Events()
	ss.api = s.API()
	ss.data = s.Data()

	ss.events.Emit(SystemEventsSelector, APIDInitializedEvent)
}

func RegisterPlugin(plugin PluginInitFunc) {
	pluginInitFuncs = append(pluginInitFuncs, plugin)
}

func RegisterPostPlugin(plugin PostPluginInitFunc) {
	postpluginInitFuncs = append(postpluginInitFuncs, plugin)
}

func InitializePlugins() {
	log := Log()
	log.Debugf("Initializing plugins...")
	for _, p := range pluginInitFuncs {
		err := p(services)
		if err != nil {
			log.Panicf("Error initializing plugin: %s", err)
		}
	}
	Events().Emit(SystemEventsSelector, PluginsInitializedEvent)
	log.Debugf("done initializing plugins")
}

func PostInitializePlugins() {

	log := Log()
	log.Debugf("Post Initializing plugins...")
	for _, p := range postpluginInitFuncs {
		err := p(services)
		if err != nil {
			log.Panicf("Error initializing plugin: %s", err)
		}
	}
	Events().Emit(SystemEventsSelector, PostPluginsInitializedEvent)
	log.Debugf("Done Post initializing plugins")
}

func AllServices() Services {
	return services
}

func Log() LogService {
	return services.Log()
}

func API() APIService {
	return services.API()
}

func Config() ConfigService {
	return services.Config()
}

func Data() DataService {
	return services.Data()
}

func Events() EventsService {
	return services.Events()
}

type servicesSet struct {
	config ConfigService
	log    LogService
	api    APIService
	data   DataService
	events EventsService
}

func (s *servicesSet) API() APIService {
	return s.api
}

func (s *servicesSet) Config() ConfigService {
	return s.config
}

func (s *servicesSet) Data() DataService {
	return s.data
}

func (s *servicesSet) Events() EventsService {
	return s.events
}

func (s *servicesSet) Log() LogService {
	return s.log
}

type systemEvent struct {
	description string
}
