
# Determine the tag for the docker image
ifeq ($(strip $(DOCKER_TAG)),)
	DOCKER_TAG=local
endif

train:
	PYTHONPATH=$(PYTHONPATH):$(PWD) pipenv run train

run:
	pipenv run start

build:
	docker build -t predict:$(DOCKER_TAG) .

clean:
	rm -rf output/

local-all: | train run

