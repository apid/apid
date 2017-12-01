#!/bin/bash
mkdir certs
PAR="${1:-2048}"
echo $PAR
openssl req -newkey rsa:$PAR -nodes -keyout certs/key.pem -x509 -days 365 -out certs/certificate.pem
echo "Please put the following into you apid_config.yaml"
echo "api_tls_key: certs/key.pem"
echo "api_tls_cert: certs/certificate.pem"
