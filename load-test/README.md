# Docker build for load testing

## Build

    make docker

## Run
Open `Makefile` and change all the values `CLOUD_IP=192.168.99.100` to your docker host(In my example that is `docker-machine ip`. For you it might be `localhost`)

    make dockerloadtest
    watch -n 1 docker logs --tail=10 apid-lt

## Start the load
+ Install https://artillery.io/
+ open `artillery/deployment.yaml` update the `target` to point to your `CLOUD_IP`:9000
+ `artillery run artillery/deployment.yaml `
+ generate report `artillery report [file.json]`

## Start the test on AWS

We use https://github.com/f1erro/spicer
Docker commands file for spicer can be found in `docker-commands` file in this folder