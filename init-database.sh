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

    
    CREATE SEQUENCE card_seq START 1;

    CREATE TABLE cards (
        id serial PRIMARY KEY,
        issuer_id integer NOT NULL,
        owner_id integer NOT NULL,
        active boolean NOT NULL,
        usages integer NOT NULL,
        capacity integer NOT NULL
    );

    INSERT INTO cards (id, issuer_id, owner_id, active, usages, capacity)
    VALUES (nextval('card_seq'), 1, 1, TRUE, 10, 15);

    INSERT INTO cards (id, issuer_id, owner_id, active, usages, capacity)
    VALUES (nextval('card_seq'), 1, 2, TRUE, 2, 10);

    CREATE SEQUENCE event_seq START 1;

    CREATE TABLE events (
        id serial PRIMARY KEY,
        card_id integer NOT NULL,
        timestamp TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
        quantity integer NOT NULL
    );

EOSQL