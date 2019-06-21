#!/bin/bash

docker-compose up -d
npm install
sleep 10 # wait until mountebank is ready
npm test
