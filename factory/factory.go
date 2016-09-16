package factory

import (
	"github.com/30x/apid"
	"github.com/30x/apid/api"
	"github.com/30x/apid/config"
	"github.com/30x/apid/data"
	"github.com/30x/apid/events"
	"github.com/30x/apid/logger"
)

// Don't use values directly - pass to apid.Initialize()
// eg. apid.Initialize(factory.DefaultServicesFactory())

func DefaultServicesFactory() apid.Services {
	return &defaultServices{}
}

type defaultServices struct {
}

func (d *defaultServices) API() apid.APIService {
	return api.CreateService()
}

func (d *defaultServices) Config() apid.ConfigService {
	return config.GetConfig()
}

func (d *defaultServices) Data() apid.DataService {
	return data.CreateDataService()
}

func (d *defaultServices) Events() apid.EventsService {
	return events.CreateService()
}

func (d *defaultServices) Log() apid.LogService {
	return logger.Base()
}
