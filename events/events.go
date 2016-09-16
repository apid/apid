package events

import "github.com/30x/apid"

// simple pub/sub to deliver events to listeners based on a selector string

const configChannelBufferSize = "events_buffer_size"

var log apid.LogService
var config apid.ConfigService

func CreateService() apid.EventsService {
	if log == nil {
		log = apid.Log().ForModule("events")
		config = apid.Config()
		config.SetDefault(configChannelBufferSize, 5)
	}
	return &eventManager{}
}
