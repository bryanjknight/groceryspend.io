# Server

This directory holds the following business logic:
* User administration
* Receipt Parsing
* Receipt Analytics

## Databases
### PostgreSQL
We use Postgres for all of our data storage needs

## Dev Tools
* `golang-migrate`: `brew install golang-migrate`

```
# For CICD
- name: Install golang-migrate
  run: |
    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.12.2/migrate.linux-amd64.tar.gz | tar xvz
    sudo mv migrate.linux-amd64 /usr/bin/
    which migrate
```
**NOTE**: We may bring these into the `go.mod` for easier usage in CICD

### Adding a new service with a database
1. Install dev tools above
1. run `mkdir -p services/<service-name>/db/migration`
1. run `migrate create -ext sql -dir services/<service-name>/db/migration -seq init_schema`
1. Add your create schema in init_schema.up.sql
1. Drop your tables in init_schema.down.sql
1. Create your your database (e.g. `CREATE DATABASE` or `createdb`)
1. Run `migrate -path services/<service-name>/db/migration -database "postgresql://postgres:example@localhost:5432/<service-db>?sslmode=disable" -verbose up`

### Adding a change to a database
1. Create a `<number>_name.down.sql and <number>_name.up.sql` in the db folder of the service


### Adminer
RDBMS explorer. Connect to `localhost:18080` with the following information:

| Parameter | Value |
| --------- | ------|
| Hostname  | `postgres` (the hostname in the docker network) |
| Username | `postgres` (the default super user) |
| Password | `example` |
| Database | `postgres` | 

### RabbitMQ admin console
Used for monitoring the local RabbitMQ instance. Connect to `localhost:15672` and login with `guest/guest`