package main

import (
	// import plugins to ensure they're bound into the executable
	_ "github.com/30x/apidApigeeSync"
	_ "github.com/30x/apidVerifyAPIKey"
	_ "github.com/30x/apidGatewayDeploy"

	// other imports
	"github.com/30x/apid"
	"github.com/30x/apid/factory"
	"flag"
	"os"
)

func main() {
	configFlag := flag.String("config", "", "path to the yaml config file [./apid_config.yaml]")

	flag.Parse()

	configFile := *configFlag
	if configFile != "" {
		os.Setenv("APID_CONFIG_FILE", configFile)
	}

	apid.Initialize(factory.DefaultServicesFactory())

	log := apid.Log()
	log.Debug("initializing...")

	apid.InitializePlugins()

	// start client API listener
	log.Debug("listening...")

	api := apid.API()
	err := api.Listen() // doesn't return if no error
	log.Fatal("Is APID already running?", err)
}
