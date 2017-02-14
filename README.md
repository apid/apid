# apid

apid is a container for publishing APIs that provides core services to its plugins including configuration, 
API publishing, data access, and a local pub/sub event system.

## To build and run

    glide install --strip-vendor
    go build
    ./apid

For command line options:

    ./apid -help

## Configuration

Configuration can be done via yaml file or environment variables. Keys are case-insensitive. 
By default, apid will look for a file called apid_config.yaml in the current working directory.
See apid_config_sample.yaml in this directory for an example file.

### Defaults

    api_port: 9000
    data_path: /var/tmp
    log_level: debug    # valid values: Debug, Info, Warning, Error, Fatal, Panic 
 
### Environment variable overrides

Config will pick up env vars automatically. Use "apid_" as a prefix for settings. For example, for 
apid's "log_level" configuration setting, set env var "apid_log_level". 

## Helpful Hints

* Use `export APID_DATA_TRACE_LOG_LEVEL=debug` to see DB Tracing


## Components
 
### Core

* [ApigeeSync](https://github.com/30x/apidApigeeSync)
 
### Plugins

* [ApigeeSync](https://github.com/30x/apidApigeeSync)
* [VerifyAPIKey](https://github.com/30x/apidVerifyApiKey)
* [GatewayDeploy](https://github.com/30x/apidGatewayDeploy)
* [ApigeeAnalytics](https://github.com/30x/apidAnalytics)

To change plugins list, edit main.go and update glide.yaml.

### Glide Build

In order to have a reproducible build, we are following the following rules:

* All glide.yaml packages must specify the version. 
* Any 3rd-party libraries used by apid-core or any plugin must be checked into the vendor directory here.
* Check in glide.lock and only use `glide up` when updating library versions
* Each plugin's glide.yaml versions should match the libraries versions in apid
* Use an empty $GOPATH to ensure a clean build


Note: Use `glide install --strip-vendor` to keep vendor directory size to a minimum

#### Important 

apid-core and all plugins must only rely on libraries including by this apid module.
If there are additional libraries that are needed, they must be approved and added to glide.yaml
and checked in to the vendor directory. 
