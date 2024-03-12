#!/bin/bash

set -o errexit
set -o pipefail
set -o nounset

# Cleanup resources
function cleanup {
  set +e
  ./stop.sh
  cd ..
  rm -rf output
  set -e
}
trap cleanup EXIT

function waitAlfrescoReady {
  echo "Starting Alfresco ..."
  bash -c 'while [[ "$(curl -s -o /dev/null -w ''%{http_code}'' http://localhost/alfresco/s/api/server)" != "200" ]]; do sleep 5; done'
  echo "Alfresco started successfully!"
}

cd ..

# Alfresco Community
go run main.go create -v 23.2

cd output

./start.sh

waitAlfrescoReady

cleanup