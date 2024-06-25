#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE USER docker;   

    CREATE DATABASE loyality_cards;
    GRANT ALL PRIVILEGES ON DATABASE loyality_cards TO docker;

    \c loyality_cards;

    CREATE SEQUENCE client_seq START 1;

    CREATE TABLE clients (
        id serial PRIMARY KEY,
        name VARCHAR(50) UNIQUE NOT NULL,
        email VARCHAR(50) UNIQUE NOT NULL
    );

    INSERT INTO clients(id, name, email)
    VALUES (nextval('client_seq'), 'foo', 'foo@gmail.com');

    INSERT INTO clients(id, name, email)
    VALUES (nextval('client_seq'), 'bar', 'bar@wp.pl');

    INSERT INTO clients(id, name, email)
    VALUES (nextval('client_seq'), 'grogu', 'grogu@onet.pl');


    CREATE SEQUENCE issuer_seq START 1;

    CREATE TABLE issuers (
        id serial PRIMARY KEY,
        name VARCHAR(50) UNIQUE NOT NULL
    );

    INSERT INTO issuers(id, name)
    VALUES (nextval('issuer_seq'), 'Abc');

    INSERT INTO issuers(id, name)
    VALUES (nextval('issuer_seq'), 'Def');

    INSERT INTO issuers(id, name)
    VALUES (nextval('issuer_seq'), 'Ghi');
EOSQL