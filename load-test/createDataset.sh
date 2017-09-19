#!/bin/bash

ORG=eaptest5
MGMT_HOST=https://api.e2e.apigee.net
MGMT_USER_ID=alexkhimich+edgex-apid-tests@google.com
MGMT_USER_PWD=Google123
NO_OF_DEVELOPERS_PER_PRODUCT=1000
DEV_PREFIX=apid_test
NO_OF_APPS_PER_DEV=1
NO_OF_PRODUCTS=15

CREATE_PRODUCT_BODY_TEMPLATE='
<ApiProduct name="test_product_name">
    <DisplayName>test_product_name</DisplayName>
    <ApprovalType>auto</ApprovalType>
    <ApiResources>
        <ApiResource>/**</ApiResource>
        <ApiResource>/</ApiResource>
    </ApiResources>
    <Environments>
        <Environment>test</Environment>
    </Environments>
    <Attributes>
      <Attribute>
         <Name>Company</Name>
         <Value>Apigee</Value>
      </Attribute>
    </Attributes>
</ApiProduct>
'


for (( i=0; i<NO_OF_PRODUCTS ; i++ ))
do
    echo -n "creating product : $i "
    PRODUCT_NAME="test_product_"${i}
    echo ${CREATE_PRODUCT_BODY_TEMPLATE} | sed "s/test_product_name/"${PRODUCT_NAME}"/g" > ./temppostbody
    curl -X POST ${MGMT_HOST}/v1/o/${ORG}/apiproducts -u ${MGMT_USER_ID}:${MGMT_USER_PWD} -H "Content-Type: application/xml" --data "@./temppostbody"
    sleep 1
    for (( c=0; c<$NO_OF_DEVELOPERS_PER_PRODUCT ; c++ ))
    do
	    echo -n "creating developer : $c "
	    DEVELOPER_NAME='edgex_load_test_dev_'${PRODUCT_NAME}'_'${c}
	    DEVELOPER_EMAIL=${DEVELOPER_NAME}'@google.com'
	    curl -X POST ${MGMT_HOST}/v1/o/${ORG}/developers -u ${MGMT_USER_ID}:${MGMT_USER_PWD} -H "Content-Type: application/json" \
	    --data '{"email": "'${DEVELOPER_EMAIL}'", "firstName": "'${DEVELOPER_NAME}'", "lastName": "load-test", "userName": "'${DEVELOPER_NAME}'", "attributes": [{ "name": "region", "value": "north" }]}'
	    sleep 0.2
	    for (( a=0; a<$NO_OF_APPS_PER_DEV ; a++ ))
	    do
	        echo -n "creating app for developer $c : $a "
	        curl -X POST ${MGMT_HOST}/v1/o/${ORG}/developers/${DEVELOPER_EMAIL}/apps -u ${MGMT_USER_ID}:${MGMT_USER_PWD} -H "Content-Type: application/json" \
	        --data '{"accessType" : "read","attributes" : [ { "name" : "Language", "value" : "java"  } ],  "apiProducts": [ "'${PRODUCT_NAME}'"],  "callbackUrl" : "www.apigee.com",  "name" : "'${DEVELOPER_NAME}'_'${a}'"}'
        done
    done
done