# Server

This directory holds the following business logic:
* User administration
* Receipt Parsing
* Receipt Analytics

## Databases
### MongoDB
We use MongoDB as a document store for the original receipt requests. This could change over time to just a file store or a different document DB.

### PostgreSQL
We use Postgres for all of our usual RDBMS needs. We could potentially use this for documents as well; however, blob storage in RDBMS is traditionally a no-no. This requires more research

### Redis
We use Redis for session management

## Dev Tools

### Mongo Express
Mongo DB explorer. Connect to `localhost:18081` for access, default creds are `admin:pass`

### Adminer
RDBMS explorer. Connect to `localhost:18080` with the following information:

| Parameter | Value |
| --------- | ------|
| Hostname  | `postgres` (the hostname in the docker network) |
| Username | `postgres` (the default super user) |
| Password | `example` |
| Database | `postgres` | 


### Redis Commander
Redis explorer. Connect to `localhost:218081`, default creds are `root:qwerty`