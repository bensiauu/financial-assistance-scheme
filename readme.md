# Financial Assistance Scheme API

This project is a Go-based API for managing financial assistance schemes, applicants, and applications. The API uses PostgreSQL as the database and secures endpoints using JWT-based authentication.

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Setup and Run the Development Environment](#setup-and-run-the-development-environment)
   - [Running with Docker Compose](#running-with-docker-compose)
   - [Local Deployment](#local-deployment)
3. [Testing](#testing)
4. [API Documentation]()
5. [Troubleshooting](#troubleshooting)

## Prerequisites

Before setting up the development environment, ensure you have the following installed on your machine:

- **Docker**: [Install Docker](https://docs.docker.com/get-docker/)
- **Docker Compose**: [Install Docker Compose](https://docs.docker.com/compose/install/)
- **Go (Golang)**: [Install Go](https://golang.org/doc/install) (Optional, for local development)
- **Postman**: [Install Postman](https://www.postman.com/downloads/) (for API testing and documentation)

## Setup and Run the Development Environment

### Running with Docker Compose

This project includes a script `generate_password_and_run.sh` that:

1. Generates a random password for the PostgreSQL database.
2. Updates the `docker-compose.yml` file with the generated password.
3. Exports the necessary environment variables.
4. Runs the application using Docker Compose.

To run the application using Docker Compose:

```bash
chmod +x generate_password_and_run.sh
./generate_password_and_run.sh
```

This script will:

- Generate a random password for the PostgreSQL database.
- Start the database and the application services.
- Apply any pending migrations to the database.
- Expose the API on `localhost:8080`.

You can stop the containers using:

```bash
docker-compose down
```

### Local Deployment

For local deployment (without Docker), follow these steps:

1. **Clone the Repository**:

   ```bash
   git clone https://github.com/bensiauu/financial-assistance-scheme.git
   cd financial-assistance-scheme
   ```

2. **Install Dependencies**:
   Ensure you have Go installed, then install the dependencies:

   ```bash
   go mod tidy
   ```

3. **Setup PostgreSQL**:
   Ensure you have a local PostgreSQL instance running. You can set it up as follows:

   - **Install PostgreSQL**: Follow the instructions on [PostgreSQL's website](https://www.postgresql.org/download/).
   - **Create a Database and User**:

     ```sql
     CREATE USER govtech WITH ENCRYPTED PASSWORD 'password123';
     CREATE DATABASE financial_assistance;
     ALTER DATABASE financial_assistance OWNER TO govtech;

     <!-- connect to financial_assistance database as govtech -->
     CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
     CREATE EXTENSION IF NOT EXISTS "pgcrypto";


     ```

4. **Configure Environment Variables**:
   Set the required environment variables:

   ```bash
   export DB_USER=govtech
   export DB_PASSWORD=password123
   export DB_NAME=financial_assistance
   export DB_HOST=localhost
   export DB_PORT=5432
   ```

5. **Run the Application**:
   Start the application:

   ```bash
   go run cmd/api/main.go
   ```

   The API will now be running at `http://localhost:8080`.

## Testing

### Running Tests with Docker

To run your tests inside Docker using Docker Compose:

```bash
docker-compose run --rm app go test ./...
```

### Running Tests Locally

To run the tests locally:

```bash
go test ./...
```

Ensure the PostgreSQL test database is running with the correct environment variables for testing.

## API Documentation

API documentation can be found [here](https://documenter.getpostman.com/view/8685199/2sAXjNYB7C).

## Troubleshooting

### Common Issues

1. **Database Connection Issues**:

   - Ensure that the `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, and `DB_NAME` environment variables are set correctly.
   - Ensure that the PostgreSQL service is running and accessible.

2. **Migrations Not Applied**:

   - Check if the migrations have been run by using `go run cmd/migrate/main.go` or by checking the `migrations` table in the database.

3. **Token Expiry or Invalid Token**:
   - Ensure that the JWT secret and token expiration time are set properly in the code and environment.

### Logs

To check the logs for the running application, use:

```bash
docker-compose logs
```
