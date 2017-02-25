#!/bin/bash

cd /demo

if [ -z "$CLOUD_IP" ];  then
  echo "No replace"
else
  sed -i "s/localhost/${CLOUD_IP}/" apid_config.yaml
fi

echo "----- Apid config being loaded -----"
cat apid_config.yaml
echo "--------- End apid config ----------"
./apid -clean
