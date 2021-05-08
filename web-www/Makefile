
# Determine the tag for the docker image
ifeq ($(strip $(DOCKER_TAG)),)
	DOCKER_TAG=local
endif

build:
	hugo --minify

docker:
	docker build -t web-www:$(DOCKER_TAG) .

