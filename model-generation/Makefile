
# Determine the tag for the docker image
ifeq ($(strip $(DOCKER_TAG)),)
	DOCKER_TAG=local
endif

train:
	PYTHONPATH=$(PYTHONPATH):$(PWD) pipenv run train

run:
	CAT_MODEL=output/categorize/model/categorize_tuple.pkl \
	CAT_YAML=data/canonical/categories.yml \
	pipenv run start

build:
	docker build -t predict:$(DOCKER_TAG) .

clean:
	rm -rf output/

local-all: | train run

