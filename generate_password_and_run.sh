#!/bin/bash

# Check if the volume exists
VOLUME_NAME="financial-assistance-scheme_postgres_data"
VOLUME_EXISTS=$(docker volume ls -q | grep -w $VOLUME_NAME)

if [ -z "$VOLUME_EXISTS" ]; then
  # Volume doesn't exist, generate a new password
  POSTGRES_PASSWORD=$(openssl rand -base64 16)

  # Create a fresh .env file with the generated password
  echo "POSTGRES_PASSWORD=${POSTGRES_PASSWORD}" >.env
  echo "DB_PASSWORD=${POSTGRES_PASSWORD}" >>.env

  echo "Generated new password and created a fresh .env file."
else
  echo "Volume already exists, skipping password generation."
fi

# Remove containers without removing volumes
docker-compose down

# Rebuild and bring up the services
docker-compose up --build
