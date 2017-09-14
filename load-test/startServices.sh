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

cd /demo

BASE_DB_PATH="/demo/data/sqlite/common/base"
DATA_DB_PATH="/demo/data/sqlite/common/1:1:"

mkdir -p $BASE_DB_PATH
mkdir -p $DATA_DB_PATH

cat init_base_db.sql | sqlite3 ${BASE_DB_PATH}/sqlite3
cat init_data_db.sql | sqlite3 ${DATA_DB_PATH}/sqlite3

echo "APID table:"
echo "select * from APID;" | sqlite3 ${BASE_DB_PATH}/sqlite3

echo "----- Apid config being loaded -----"
cat apid_config.yaml
echo "--------- End apid config ----------"
./apid
