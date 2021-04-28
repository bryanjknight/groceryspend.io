#!/bin/bash
set -e

# create databases
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE USER users;
    CREATE DATABASE users;
    GRANT ALL PRIVILEGES ON DATABASE users TO users;

    CREATE USER receipts;
    CREATE DATABASE receipts;
    GRANT ALL PRIVILEGES ON DATABASE receipts TO receipts;
EOSQL

# install uuid extesnsion on users
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "users" <<-EOSQL
    CREATE EXTENSION "uuid-ossp" 
EOSQL

# install uuid extesnsion on receipts
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "receipts" <<-EOSQL
    CREATE EXTENSION "uuid-ossp" 
EOSQL