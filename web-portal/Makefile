
# Determine the tag for the docker image
ifeq ($(strip $(DOCKER_TAG)),)
	DOCKER_TAG=local
endif

build:
	docker build -t web-portal:$(DOCKER_TAG) .

clean:
	rm -rf dist/

push:
	docker tag web-portal:$(DOCKER_TAG) groceryspend/web-portal:$(DOCKER_TAG)
	docker push groceryspend/web-portal:$(DOCKER_TAG)

