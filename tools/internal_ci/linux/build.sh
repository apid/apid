#!/bin/bash

BUILDROOT=${BUILDROOT:-git/apid}
export BUILDROOT

# Make a temporary GOPATH to build in
gobase=`mktemp -d`
GOPATH=${gobase}
export GOPATH

go get github.com/Masterminds/glide

base=${gobase}/src/github.com/30x/apid
mkdir -p ${base}
(cd ${BUILDROOT}; tar cf - .) | (cd ${base}; tar xf -)

set +x

(cd ${base}; ${GOPATH}/bin/glide install)
(cd ${base}; go build -o apid ./cmd/apid)
buildResult=$?

cp ${base}/apid .

rm -rf ${gobase}

exit ${buildResult}
