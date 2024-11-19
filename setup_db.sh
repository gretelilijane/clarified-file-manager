#!/bin/bash

# Load environment variables from .env
if [ -f .env ]; then
    export $(grep -v '^#' .env | xargs)
else
    echo ".env file not found. Aborting."
    exit 1
fi

# Check if required variables are set
if [[ -z "$POSTGRES_USER" || -z "$POSTGRES_PASSWORD" || -z "$POSTGRES_DB" || -z "$DATABASE_HOST" || -z "$DATABASE_PORT" ]]; then
    echo "Missing required environment variables. Ensure .env file is configured properly."
    exit 1
fi

# Create the database and run init.sql
echo "Setting up the database..."

# Use PSQL commands
PGPASSWORD=$POSTGRES_PASSWORD psql -h $DATABASE_HOST -p $DATABASE_PORT -U $POSTGRES_USER <<EOSQL
-- Create the database
CREATE DATABASE $POSTGRES_DB;
EOSQL

if [ $? -ne 0 ]; then
    echo "Failed to create database. Aborting."
    exit 1
fi

echo "Database $POSTGRES_DB created successfully."

# Initialize schema using init.sql
PGPASSWORD=$POSTGRES_PASSWORD psql -h $DATABASE_HOST -p $DATABASE_PORT -U $POSTGRES_USER -d $POSTGRES_DB -f ./db/init.sql

if [ $? -ne 0 ]; then
    echo "Failed to initialize database schema. Aborting."
    exit 1
fi

echo "Database schema initialized successfully."