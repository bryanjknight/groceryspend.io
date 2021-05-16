#!/bin/bash
set -e

# create databases
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE USER users PASSWORD 'users-password';
    CREATE DATABASE usersdb;
    GRANT ALL PRIVILEGES ON DATABASE usersdb TO users;

    CREATE USER receipts PASSWORD 'receipts-password';
    CREATE DATABASE receiptsdb;
    GRANT ALL PRIVILEGES ON DATABASE receiptsdb TO receipts;
EOSQL

# install uuid extesnsion on users
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "usersdb" <<-EOSQL
    CREATE EXTENSION "uuid-ossp" 
EOSQL

# install uuid extesnsion on receipts
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "receiptsdb" <<-EOSQL
    CREATE EXTENSION "uuid-ossp" 
EOSQL