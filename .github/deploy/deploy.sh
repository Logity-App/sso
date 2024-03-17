#!/bin/bash
IMAGE_NAME=$DOCKER_USERNAME/$DOCKER_REPOSITORY:$DOCKER_TAG
CONTAINER_NAME=$DOCKER_REPOSITORY_$DOCKER_TAG

docker login -u $DOCKER_USERNAME -p $DOCKER_PASSWORD
docker-compose -f .github/deploy/docker-compose.${ENVIRONMENT}.yml down --remove-orphans
docker-compose -f .github/deploy/docker-compose.${ENVIRONMENT}.yml pull
docker-compose -f .github/deploy/docker-compose.${ENVIRONMENT}.yml up -d