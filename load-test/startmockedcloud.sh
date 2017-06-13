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

export BUNDLE_URI="http://$CLOUD_IP:$APID_API_PORT/bundles/1"
./mockServer -numDeps=100 -numDevs=50000 -addDevEach=3s -upDevEach=1s -upDepEach=3s -bundleURI=$BUNDLE_URI