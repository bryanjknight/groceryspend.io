# Server

This directory holds the following business logic:
* User administration
* Receipt Parsing
* Receipt Analytics

## Databases
### PostgreSQL
We use Postgres for all of our data storage needs

## Dev Tools
* `sqlc`: `brew install sqlc`
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

### Adding queries via `sqlc`
1. Install `sqlc` as described above
1. Update service in `sqlc.yaml`. Below is an example:
```yaml
  - name: "users"
    path: "./services/users"
    queries: "./services/users/db/queries.sql"
    schema: "./services/users/db/migration/"
    engine: "postgresql"
    emit_prepared_queries: true
    emit_interface: false
    emit_exact_table_names: false
    emit_empty_slices: false
    emit_json_tags: true
    json_tags_case_style: "camel"
```
1. Run `sqlc compile` to check sql types
1. Run `sqlc generate` to generate Go files

### Adminer
RDBMS explorer. Connect to `localhost:18080` with the following information:

| Parameter | Value |
| --------- | ------|
| Hostname  | `postgres` (the hostname in the docker network) |
| Username | `postgres` (the default super user) |
| Password | `example` |
| Database | `postgres` | 

