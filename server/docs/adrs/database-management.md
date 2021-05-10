Database Management
===

## Status
Accepted on 2021-05-09

## Context
The goal of this ADR is to pick a pattern for managing our database (e.g. migrations) as well as querying it

## Options
### Gorm
Pros:
* ORM built in
* Auto-migrate schemas
* Already implemented as part of POC

Cons:
* No clear way to do specific migrations (e.g. create null column, add data to new column, make column not null)


### SQLC and Go-Migrate
Pros:
* Go-migrate has a similar framework to django, which is nice, supports idepotency and rollback
* SQLC leverages code generation to crate typed intefaces

Cons:
* Requires some rewrite of the postgres db repository
* Associated objects kind of hard to implement

### SQLX and Go-Migrate
Pros:
* Go-migrate has a similar framework to django, which is nice, supports idepotency and rollback
* SQLX leverages existing database/sql library with `StructScan`

Cons:
* Requires some rewrite of the postgres db repository (less that sqlc)
* Associated objects by hand

## Decision
* We will implement via SQLX and Go-Migrate

## Consequences
* Need to think about how to manage the number of migration scripts

## Compliance
* TBD

## Notes
