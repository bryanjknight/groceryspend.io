
# Determine the tag for the docker image
ifeq ($(strip $(DOCKER_TAG)),)
	DOCKER_TAG=local
endif

lint:
	golint -set_exit_status ./...

build:
	docker build -t server:$(DOCKER_TAG) .

unit-test:
	go test --cover -v ./...

run:
	docker compose up -d
	go run main.go

clean:
	rm -rf docker_data/pgdata/

all: lint unit-test build