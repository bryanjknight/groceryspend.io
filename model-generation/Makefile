
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
	docker tag predict:$(DOCKER_TAG) groceryspend/predict:$(DOCKER_TAG)

clean:
	rm -rf output/

push:
	docker push groceryspend/predict:$(DOCKER_TAG)


local-all: | train run

