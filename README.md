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

    api_listener: 127.0.0.1:9000  # api listener will bind to specified host and port
    data_path: /var/tmp           # all data stored on disk will be under this path 
    log_level: debug              # valid values: Debug, Info, Warning, Error, Fatal, Panic 
 
### Environment variable overrides

Config will pick up env vars automatically. Use "apid_" as a prefix for settings. For example, for 
apid's "log_level" configuration setting, set env var "apid_log_level". 

## Helpful Hints

* Use `export APID_DATA_TRACE_LOG_LEVEL=debug` to see DB Tracing


## Components
 
### Core

* [ApigeeSync](https://github.com/apid/apidApigeeSync)
 
### Plugins

* [ApigeeSync](https://github.com/apid/apidApigeeSync)
* [VerifyAPIKey](https://github.com/apid/apidVerifyApiKey)
* [GatewayDeploy](https://github.com/apid/apidGatewayDeploy)
* [ApigeeAnalytics](https://github.com/apid/apidAnalytics)

To change plugins list, edit main.go add to glide.yaml and follow the release process below.

### Release process

To update the build dependencies and release, follow this process:

1. Update `glide.yaml` to the correct versions
2. `rm glide.lock`
3. Run `glide up`
4. `go build` and verify the result
5. Commit (should be glide.yaml, glide.lock, and vendor changes)
6. Add a Git label on the commit for the version of the release (use semver, ie. "1.2.3") 
7. Push commit and label to Github to cause Travis to create a release and attach binaries for OS X and Linux
8. Once the release has been created, you may edit it on Github to add release notes  

#### Notes on release process

In order to have a reproducible build, we have the following rules:

* All `glide.yaml` packages must specify a version. 
* Any 3rd-party libraries used by apid-core or used by any plugin must be checked into the vendor dir.
* The `glide.lock` file is checked in and `glide up` is only used to update library versions
* Use an empty $GOPATH (aside from apid itself) to ensure a clean build
* apid-core and all plugins must only rely on libraries including by this apid module.
  If there are additional libraries that are needed, they must be approved, added to glide.yaml,
  and checked in to the vendor directory. 
