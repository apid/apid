package main

import (
	// import plugins to ensure they're bound into the executable
	_ "github.com/30x/apidAnalytics"
	_ "github.com/30x/apidApigeeSync"
	_ "github.com/30x/apidGatewayDeploy"
	_ "github.com/30x/apidVerifyAPIKey"

	// other imports
	"flag"
	"fmt"
	"os"
	"reflect"

	"github.com/30x/apid-core"
	"github.com/30x/apid-core/factory"
)

func main() {
	// clean exit messages w/o stack track during initialization
	defer func() {
		if r := recover(); r != nil {
			// if coming from logrus, it's already been printed. otherwise...
			if reflect.TypeOf(r).String() != "*logrus.Entry" {
				fmt.Println(r)
			}
			os.Exit(1)
		}
	}()

	f := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	configFlag := f.String("config", "", "path to the yaml config file [./apid_config.yaml]")
	cleanFlag := f.Bool("clean", false, "start clean, deletes all existing data from local_storage_path")

	f.Parse(os.Args[1:])

	configFile := *configFlag
	if configFile != "" {
		os.Setenv("APID_CONFIG_FILE", configFile)
	}

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

	log.Infof("Wait for plugins to gracefully shutdown")
	apid.ShutdownPluginsAndWait()
	log.Infof("Apid graceful shutdown succeeded")
}
