DOCKER_IMAGE_NAME=redirect-service
DOCKER_CONTAINER_NAME=redirect-service-container

.PHONY: build run

build:
	docker build -t $(DOCKER_IMAGE_NAME) .

run:
	docker stop $(DOCKER_CONTAINER_NAME) || true
	docker rm $(DOCKER_CONTAINER_NAME) || true
	docker run --name $(DOCKER_CONTAINER_NAME) -p 80:80 $(DOCKER_IMAGE_NAME)
