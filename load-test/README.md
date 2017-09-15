# Docker build for load testing

## Build
    make docker

## Run
    make dockerloadtest
    watch -n 1 docker logs --tail=10 apid-lt

## Clean
	make clean

## Install Artillery
+ Install https://artillery.io/
### For Googlers using Gubuntu:
+ Please follow go/nodejs-install to intall nodejs locally.
+ Install https://artillery.io/

## Start the load
+ open `artillery/configurations.yaml` update the `target` to point to your `CLOUD_IP`:9000
+ `artillery run artillery/configurations.yaml `
+ generate report `artillery report [file.json]`

