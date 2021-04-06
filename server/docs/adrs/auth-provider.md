Auth and Authz providers
===

## Status
Accepted on 2021-04-06

## Context
The goal of this ADR is to pick a design paradigm around authentication and authorization.

## Options
### Auth0
Pros:
* Session management, auth providers all out of the box
* Free package gives 7k MAU with 2 socail connections
* JWT is managed by Auth0, can reduce the number of requests to Auth0 to verify validity

Cons:
* Killing an active session is not possible since the JWT dictates when the session expires


### Our own auth/authz
Pros:
* We can kill sessions by removing session info from Redis
* Auth providers are well known and libraries exist
* No MAU costs; however, there's an increase cost of upfront implementation

Cons:
* Easy to make mistake, creating a vulnerability

## Decision
* We will implement via Auth0, but keep the middleware generic so that if we want to roll our own middleware, it's not a massive lift

## Consequences
* We need to keep an eye on MAU and bad sessions

## Compliance
* TBD

## Notes
