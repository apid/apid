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

(cd ${base}; ${GOPATH}/bin/glide install)
(cd ${base}; go test ./api ./config ./events ./factory ./logger)
testResult=$?

if [ $testResult -eq 0 ]
then 
  echo "Building apid binary"
  (cd ${base}; go build ./cmd/apid)
  testResult=$?
fi

rm -rf ${gobase}

exit ${testResult}
