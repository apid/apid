#!/bin/bash

export BUNDLE_URI="http://$CLOUD_IP:$APID_API_PORT/bundles/1"
./mockServer -numDeps=100 -numDevs=50000 -addDevEach=3s -upDevEach=1s -upDepEach=3s -bundleURI=$BUNDLE_URI