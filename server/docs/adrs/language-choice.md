Language Chioce
===

## Status
Accepted on 2021-03-25

## Context
The goal of this ADR is to pick a language for our server

## Options
### Golang
Pros:
* Performance is the best of the three
* Strong typing
* Very small binary

Cons:
* Not as experienced with the language, so there's a learning curve
* Ecosystem are still a work in progress compared to pip and npm

### Python
Pros:
* Previous project was with Python/Django
* Django could be a force multiplier to build out the features

Cons:
* Typing is a bolt-on at best
* Poor performance and heavy weight docker images

### NodeJS w/Typescript
Pros:
* Very experienced with this stack
* Robust ecosystem

Cons:
* Not compiled, so slower performance
* Somewhere between Golang and Python in terms of docker size

## Decision
* We will implement via Golang

## Consequences
* Need to think about how to manage gaps in missing libraries

## Compliance
* TBD

## Notes
