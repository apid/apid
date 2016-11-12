# cmd/apid 

## Core plugins
* [ApigeeSync](https://github.com/30x/apidApigeeSync)
* [VerifyAPIKey](https://github.com/30x/apidVerifyApiKey)
* [GatewayDeploy](https://github.com/30x/apidGatewayDeploy)

To change plugins list, edit main.go and update glide.yaml.

## Build and execute

    glide install
    go build
    ./apid

For options:

    ./apid -help
