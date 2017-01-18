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
		err := os.RemoveAll(localStorage)
		if err != nil {
			log.Panic("Failed to clean data directory: %v", err)
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
