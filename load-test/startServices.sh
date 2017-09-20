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

#cat init_base_db.sql | sqlite3 ${BASE_DB_PATH}/sqlite3
#cat init_data_db.sql | sqlite3 ${DATA_DB_PATH}/sqlite3
mv base_sqlite3 ${BASE_DB_PATH}/sqlite3
mv data_sqlite3 ${DATA_DB_PATH}/sqlite3
CREATE_INDEX='CREATE INDEX app_company_id on KMS_APP (company_id);
CREATE INDEX kms_app_id on KMS_APP_CREDENTIAL (app_id);
CREATE INDEX app_developer_id on KMS_APP (developer_id);
CREATE INDEX kms_id on KMS_APP_CREDENTIAL (id);
CREATE INDEX kms_tenant_id on KMS_APP_CREDENTIAL (tenant_id);
CREATE INDEX mp_appcred_id on KMS_APP_CREDENTIAL_APIPRODUCT_MAPPER (appcred_id);
CREATE INDEX mp_app_id on KMS_APP_CREDENTIAL_APIPRODUCT_MAPPER (app_id);
CREATE INDEX mp_apiprdt_id on KMS_APP_CREDENTIAL_APIPRODUCT_MAPPER (apiprdt_id);
CREATE INDEX org_tenant_id on KMS_ORGANIZATION (tenant_id);
CREATE INDEX org_name on KMS_ORGANIZATION (name);
CREATE INDEX ap_id on KMS_API_PRODUCT (id);
CREATE INDEX ap_tenant_id on KMS_API_PRODUCT (tenant_id);'
echo ${CREATE_INDEX} | sqlite3 ${DATA_DB_PATH}/sqlite3
echo "APID table:"
echo "select * from APID;" | sqlite3 ${BASE_DB_PATH}/sqlite3

echo "----- Apid config being loaded -----"
cat apid_config.yaml
echo "--------- End apid config ----------"
./apid
