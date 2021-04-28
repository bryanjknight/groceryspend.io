
lint:
	golint -set_exit_status ./...

build:
	# potentially use this for building a docker image
	go build -v ./...

unit-test:
	go test --cover -v ./...

run:
	docker compose up -d
	go run main.go

clean:
	rm -rf docker_data/pgdata/

all: lint build unit-test