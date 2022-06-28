#!/bin/zsh

set -e

CERTIFICATE_MATERIAL_SET="${CERTIFICATE_MATERIAL_SET:-com.amazon.certificates.dev-dsk-vanggie-2a-cc09f356.us-west-2.amazon.com-STANDARD_SSL_SERVER_INTERNAL_ENDPOINT-RSA-Chain}"
/apollo/env/envImprovement/bin/odin-get -t Principal $CERTIFICATE_MATERIAL_SET \
 | sudo tee /tmp/server-chain.crt > /dev/null

/apollo/env/envImprovement/bin/odin-get -t Credential $CERTIFICATE_MATERIAL_SET \
| sudo tee /tmp/server.key > /dev/null

# cert.pem is CA certificates from AmazonCerficiateCA.