#!/bin/bash -eu
#
# Copyright 2017 Google Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

#!/bin/bash

BUILDROOT=${BUILDROOT:-git/apid}
export BUILDROOT

# Make a temporary GOPATH to build in
gobase=`mktemp -d`
GOPATH=${gobase}
export GOPATH

go get github.com/Masterminds/glide

base=${gobase}/src/github.com/apid/apid
mkdir -p ${base}
(cd ${BUILDROOT}; tar cf - .) | (cd ${base}; tar xf -)

set +x

(cd ${base}; ${GOPATH}/bin/glide install)
(cd ${base}; go build -o apid)
buildResult=$?

cp ${base}/apid .

rm -rf ${gobase}

exit ${buildResult}
