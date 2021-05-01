
# Determine the tag for the docker image
ifeq ($(strip $(DOCKER_TAG)),)
    DOCKER_TAG=local
endif

train:
  pipenv run train

run:
	pipenv run start

build:
	docker build -t groceryspend.io/predict:$(DOCKER_TAG) .

clean:
	rm -rf output/

local-all: | train run

