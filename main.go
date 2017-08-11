// Copyright 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	// import plugins to ensure they're bound into the executable
	_ "github.com/30x/apidAnalytics"
	_ "github.com/30x/apidApigeeSync"
	_ "github.com/30x/apidGatewayConfDeploy"
	_ "github.com/30x/apidVerifyAPIKey"

	// other imports
	"flag"
	"fmt"
	"os"
	"reflect"

	"github.com/30x/apid-core"
	"github.com/30x/apid-core/factory"
	"github.com/30x/apid/version"
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
	versionFlag := f.Bool("version", false, "display the version number of apid and exits")

	f.Parse(os.Args[1:])

	if *versionFlag {
		fmt.Println("APID Version Number: " + version.VERSION_NUMBER)
		return
	}

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

	apid.InitializePlugins(version.VERSION_NUMBER)

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
