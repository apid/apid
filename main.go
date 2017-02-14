package main

import (
	// import plugins to ensure they're bound into the executable
	_ "github.com/30x/apidAnalytics"
	_ "github.com/30x/apidApigeeSync"
	_ "github.com/30x/apidGatewayDeploy"
	_ "github.com/30x/apidVerifyAPIKey"

	// other imports
	"flag"
	"os"

	"github.com/30x/apid-core"
	"github.com/30x/apid-core/factory"
)

func main() {
	configFlag := flag.String("config", "", "path to the yaml config file [./apid_config.yaml]")
	cleanFlag := flag.Bool("clean", false, "start clean, deletes all existing data from local_storage_path")

	configFile := *configFlag
	if configFile != "" {
		os.Setenv("APID_CONFIG_FILE", configFile)
	}

	flag.Parse()

	apid.Initialize(factory.DefaultServicesFactory())

	log := apid.Log()
	config := apid.Config()

	if *cleanFlag {
		localStorage := config.GetString("local_storage_path")
		log.Infof("removing existing data from: %s", localStorage)
		if err := os.RemoveAll(localStorage); err != nil {
			log.Panicf("Failed to clean data directory: %v", err)
		}
		if err := os.MkdirAll(localStorage, 0700); err != nil {
			log.Panicf("can't create local storage path %s:%v", localStorage, err)
		}
	}

	log.Debug("initializing...")

	apid.InitializePlugins()

	// start client API listener
	log.Debug("listening...")

	api := apid.API()
	err := api.Listen()
	if err != nil {
		log.Print(err)
	}
}
