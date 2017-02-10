#!/bin/bash

cd /demo

echo "----- Apid config being loaded -----"
cat apid_config.yaml
echo "--------- End apid config ----------"

APID_API_PORT=9001 ./mockServer -numDeps=100 -numDevs=50000 -addDevEach=3s -upDevEach=1s -upDepEach=3s  &
sleep 2
./apid -clean &
