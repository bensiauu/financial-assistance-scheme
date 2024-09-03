#!/bin/bash

# Generate a random 16-character password
POSTGRES_PASSWORD=$(openssl rand -base64 16)

# Create a fresh .env file with the generated password
echo "POSTGRES_PASSWORD=${POSTGRES_PASSWORD}" > .env
echo "DB_PASSWORD=${POSTGRES_PASSWORD}" >> .env

# Remove any existing containers to avoid caching issues
docker-compose down --volumes

# Run Docker Compose with the new .env file
docker-compose up --build
