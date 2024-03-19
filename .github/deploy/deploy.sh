#!/bin/bash
IMAGE_NAME=$DOCKER_USERNAME/$DOCKER_REPOSITORY:$DOCKER_TAG
CONTAINER_NAME=$DOCKER_REPOSITORY_$DOCKER_TAG

NETWORK_NAME=existing-network
if [ -z $(docker network ls --filter name=^${NETWORK_NAME}$ --format="{{ .Name }}") ] ; then
     docker network create ${NETWORK_NAME} ;
fi

docker login -u $DOCKER_USERNAME -p $DOCKER_PASSWORD
docker-compose -f .github/deploy/docker-compose.${ENVIRONMENT}.yml down --remove-orphans
docker-compose -f .github/deploy/docker-compose.${ENVIRONMENT}.yml pull
docker-compose -f .github/deploy/docker-compose.${ENVIRONMENT}.yml up $APP_NAME -d