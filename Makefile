build:
	go build -ldflags "-X main.APID_SOURCE_VERSION=`git rev-parse --abbrev-ref HEAD`-`git rev-parse HEAD`"
