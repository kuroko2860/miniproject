#!/bin/bash
set -e

# Create the database
psql -v ON_ERROR_STOP=1 --username "$kuroko" <<-EOSQL
    CREATE DATABASE miniproject;
EOSQL

# Connect to the database and create the PostGIS extension and table schema
psql -v ON_ERROR_STOP=1 --username "$kuroko" --dbname "miniproject" <<-EOSQL
    CREATE EXTENSION postgis;
    CREATE TABLE objects (
      id SERIAL PRIMARY KEY,
      object_id VARCHAR(255) NOT NULL,
      type VARCHAR(50) NOT NULL,
      color VARCHAR(50) NOT NULL,
      location GEOGRAPHY(POINT, 4326),
      status VARCHAR(50) NOT NULL,
      timestamp TIMESTAMPTZ NOT NULL
    );
    CREATE INDEX idx_objects_timestamp ON objects (timestamp);
    CREATE INDEX idx_objects_location ON objects USING GIST (location);
EOSQL
