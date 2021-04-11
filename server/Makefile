
lint:
	golint -set_exit_status ./...

build:
	# potentially use this for building a docker image
	go build -o server

unit-test:
	go test ./...

run:
	docker-compose up -d
	go run main.go

all: lint build unit-test