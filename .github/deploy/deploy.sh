#!/bin/bash
IMAGE_NAME=$DOCKER_USERNAME/$DOCKER_REPOSITORY:$DOCKER_TAG
CONTAINER_NAME=$DOCKER_REPOSITORY_$DOCKER_TAG

docker login -u $DOCKER_USERNAME -p $DOCKER_PASSWORD
docker-compose down -f .github/deploy/docker-compose.${ENVIRONMENT}.yml --remove-orphans
docker-compose pull -f .github/deploy/docker-compose.${ENVIRONMENT}.yml
docker-compose up -f .github/deploy/docker-compose.${ENVIRONMENT}.yml -d