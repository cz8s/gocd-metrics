#!/bin/bash

npm install
docker-compose up -d
while [[ "$(curl -s -o /dev/null -w ''%{http_code}'' localhost:8153/go/api/v1/health)" != "200" ]]; do 
  echo -n '.'
  sleep 5 
done
echo 'GoCD is up'

sleep 30 # gocd is up, but needs some time to discover its pipelines
npm test
