#!/bin/bash

set -e # Abort script at first error
set -u # Disallow unset variables

# Only run when not part of a pull request and on the master branch
if [ $TRAVIS_BRANCH = "master" ]
then
  docker login -u="$DOCKER_USERNAME" -p="$DOCKER_PASSWORD"
  docker tag gocd-metrics cz8s/gocd-metrics:latest
  docker push cz8s/gocd-metrics
fi
