#!/bin/bash
set -e


# we denote whether this is docker or terraform based on the fact that DigitalOcean's default
# admin user is "doadmin"

MODE="DOCKER"
SSL_REQUIRE="disable"

if [ "$POSTGRES_USER" = "doadmin" ]; then
    MODE="TERRAFORM"
    SSL_REQUIRE="require"
fi

# create users if docker (terraform *should* do this for us)
if [ "$MODE" = "DOCKER" ]; then
    psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
CREATE USER users PASSWORD 'users-password';
CREATE DATABASE usersdb;
    
CREATE USER receipts PASSWORD 'receipts-password';
CREATE DATABASE receiptsdb;
EOSQL
fi

# create databases
psql --set=sslmode="$SSL_REQUIRE" -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    GRANT ALL PRIVILEGES ON DATABASE usersdb TO users;
    GRANT ALL PRIVILEGES ON DATABASE receiptsdb TO receipts;
EOSQL

# install uuid extesnsion on users
psql --set=sslmode="$SSL_REQUIRE" -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "usersdb" <<-EOSQL
    CREATE EXTENSION "uuid-ossp" 
EOSQL

# install uuid extesnsion on receipts
psql --set=sslmode="$SSL_REQUIRE" -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "receiptsdb" <<-EOSQL
    CREATE EXTENSION "uuid-ossp" 
EOSQL