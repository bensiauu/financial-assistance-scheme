# Financial Assistance Scheme API

This project is a Go-based API for managing financial assistance schemes, applicants, and applications. The API uses PostgreSQL as the database and secures endpoints using JWT-based authentication.

## Prerequisites

Before setting up the development environment, ensure you have the following installed on your machine:

- **Docker**: [Install Docker](https://docs.docker.com/get-docker/)
- **Docker Compose**: [Install Docker Compose](https://docs.docker.com/compose/install/)
- **Go (Golang)**: [Install Go](https://golang.org/doc/install) (optional, only needed for local Go development)

## Setup and Run the Development Environment

### 1. Clone the Repository

Clone the repository to your local machine:

```bash
git clone https://github.com/yourusername/financial-assistance-scheme.git
cd financial-assistance-scheme
```

### 2. Run the Application with Docker Compose

This project includes a script `generate_password_and_run.sh` that:

1. Generates a random password for the PostgreSQL database.
2. Updates the `docker-compose.yml` file with the generated password.
3. Exports the necessary environment variables.
4. Runs the application using Docker Compose.

To run the application:

```bash
chmod +x generate_password_and_run.sh
./generate_password_and_run.sh
```

This script will:

- Generate a random password for the PostgreSQL database.
- Start the database and the application services.
- Apply any pending migrations to the database.
- Expose the API on `localhost:8080`.

### 3. Verify the Application is Running

After running the script, you should see logs indicating that the database and application services are running. You can verify the application is running by visiting:

```
http://localhost:8080
```

or by making an API request:

```bash
curl http://localhost:8080/api/some-endpoint
```

### 4. Stopping the Application

To stop the application, you can use Docker Compose:

```bash
docker-compose down
```

This will stop and remove the containers but keep the database data intact.

### 5. Cleaning Up

If you want to remove the containers and associated volumes (including the database data), use:

```bash
docker-compose down -v
```

This will stop the containers and remove the volumes.

### 2. Tests

You can see the tests that are run if you go to the github repo -> actions.

### 2. API Documentation

Documentation can be accessed [here](https://documenter.getpostman.com/view/8685199/2sAXjNYB7C)

## Troubleshooting

### Common Issues

1. **Connection Refused**: If the application can't connect to the database, ensure the database container is running and the connection details are correct.
2. **Password Authentication Failed**: Ensure the environment variables in `docker-compose.yml` match the generated password.

### Logs

Check the logs for more detailed error messages:

```bash
docker-compose logs
```

---
