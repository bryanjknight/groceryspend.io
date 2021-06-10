
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
	scripts/wait-for localhost:5432 -- echo "PostgreSQL is up"
	scripts/wait-for localhost:5672 -- echo "RabbitMQ is up"
	go run main.go

clean:
	docker compose down
	rm -rf docker_data/pgdata/
	rm -rf docker_data/rabbitmq/mnesia
	rm -f docker_data/rabbitmq/.erlang.cookie

push:
	docker tag server:$(DOCKER_TAG) groceryspend/server:$(DOCKER_TAG)
	docker push groceryspend/server:$(DOCKER_TAG)

export-ext-model:
	BROWSER_PROJECT_DIR=../browser-extension go run cmd/typescriptify/browser_ext/export.go


export-portal-model:
	PORTAL_PROJECT_DIR=../web-portal go run cmd/typescriptify/portal/export.go

migrate-payments-up:
	migrate \
		-path services/payments/db/migration \
		-database "postgres://postgres:example@localhost:5432/paymentsdb?sslmode=disable" \
		-verbose up

migrate-payments-down:
	migrate \
		-path services/payments/db/migration \
		-database "postgres://postgres:example@localhost:5432/paymentsdb?sslmode=disable" \
		-verbose down

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