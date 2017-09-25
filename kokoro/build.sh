#!/bin/bash
set -ex

BUILDROOT=${BUILDROOT:-github/apid}
export BUILDROOT

# Make a temporary GOPATH to build in
gobase=`mktemp -d`
base=${gobase}/src/github.com/apid/apid
GOPATH=${gobase}
export GOPATH

base=${GOPATH}/src/github.com/apid/apid
mkdir -p ${base}
(cd ${BUILDROOT}; tar cf - .) | (cd ${base}; tar xf -)
cd ${base}


echo "Getting glide"
go get github.com/Masterminds/glide
echo "Install dependencies for tests"
time ${GOPATH}/bin/glide up -v

go version
go build
go test $(${GOPATH}/bin/glide novendor)
