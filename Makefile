APP_NAME=logity
DOCKER_REPOSITORY=logity

docker-build:
	docker build -f build/app/Dockerfile -t $(APP_NAME):latest .
	docker tag $(APP_NAME):latest $(DOCKER_REPOSITORY)/$(APP_NAME):latest
	docker push $(DOCKER_REPOSITORY)/$(APP_NAME):latest