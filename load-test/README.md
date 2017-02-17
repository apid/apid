# Docker build for load testing

## Build

    make dockers

## Run

    docker run -d

docker will open port 9000 for api hits

## Start the load
+ Install https://artillery.io/
+ open `artillery/deployment.yaml` update the `target` to point to your dockerip:19000
+ `artillery run artillery/deployment.yaml `
+ generate report ``