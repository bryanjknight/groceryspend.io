
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

integration-test:
	go test --cover -v --tags=integration  ./...	

run:
	docker compose up -d
	go run main.go

clean:
	rm -rf docker_data/pgdata/

push:
	docker tag server:$(DOCKER_TAG) groceryspend/server:$(DOCKER_TAG)
	docker push groceryspend/server:$(DOCKER_TAG)

export-portal-model:
	PORTAL_PROJECT_DIR=../web-portal go run cmd/typescriptify/portal/export.go

migrate-receipts-up:
	migrate \
		-path services/receipts/db/migration \
		-database "postgres://postgres:example@localhost:5432/receiptsdb?sslmode=disable" \
		-verbose up

migrate-receipts-down:
	migrate \
		-path services/receipts/db/migration \
		-database "postgres://postgres:example@localhost:5432/receiptsdb?sslmode=disable" \
		-verbose down

migrate-user-up:
	migrate \
		-path services/users/db/migration \
		-database "postgres://postgres:example@localhost:5432/usersdb?sslmode=disable" \
		-verbose up

migrate-user-down:
	migrate \
		-path services/users/db/migration \
		-database "postgres://postgres:example@localhost:5432/usersdb?sslmode=disable" \
		-verbose down

all: lint unit-test build