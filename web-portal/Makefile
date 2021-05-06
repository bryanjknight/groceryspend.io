
# Determine the tag for the docker image
ifeq ($(strip $(DOCKER_TAG)),)
	DOCKER_TAG=local
endif

build:
	docker build -t web-portal:$(DOCKER_TAG) .

push:
	echo $DOCKER_PASS | docker login --username bknight1983 --password-stdin
	docker tag web-portal:$(DOCKER_TAG) groceryspend/web-portal:$(DOCKER_TAG)
	docker push groceryspend/web-portal:$(DOCKER_TAG)

