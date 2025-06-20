# Transfer System

A backend API for a transfer system, built with **Go** and **PostgreSQL**, following Clean Architecture principles

## Features

- Account management (create, retrieve)
- Transaction management (transfer)
- Clean architecture with dependency injection
- PostgreSQL support
- Logging and graceful shutdown

## Project Structure

```
adapters/         # Controllers, repositories, and web routes
cmd/              # Application entrypoint (main.go)
docs/             # Generated OpenAPI Specs
domain/           # Business logic and interface definitions (ports)
infrastructure/   # DB and external service implementations
internal/         # Internal utilities, currently only contain for test
mocks/            # Mocks for tests
pkg/              # Shared packages (logger, config, etc.)
```

## Prerequisites

- [Go 1.24+](https://go.dev/doc/install)
- [PostgreSQL](https://www.postgresql.org/)
- [Docker + Docker Compose](https://docs.docker.com/compose/)

---

## Getting Started

### 1. Clone the repository

```sh
git clone https://github.com/chud-lori/transfer-system
cd transfer-system
```

### 2. Setup environment

```sh
cp .env.example .env
```

Edit `.env` to configure:

- `DB_URL`
- `APP_PORT`
- `POSTGRES_USER`
- `POSTGRES_PASSWORD`

---

## Build and Run Options

### Option 1: Using Makefile
Make sure you have installed postgresql in your local or as docker container

create database in your postgresql

then run this

`psql -h <hostname> -p <port> -d <database_name> -U <username> -f db.sql`


if you have postgresql as docker container, copy the db.sql to container


`docker cp db.sql <container_name>:/tmp/db.sql`


and import it


`docker exec -it <container_name> psql -U postgres -d <database_name> -f /tmp/db.sql`


and build using `make`
```sh
make
```

### Option 2: Run Locally with Go

```sh
go mod download
go run ./cmd/main.go
```

### Option 3: Run with Docker Compose

```sh
docker-compose up --build
```

This will spin up the app and a PostgreSQL container using the configurations in `docker-compose.yml`.

---
## Run Tests
The test require to use postgresql database, make sure to create `.env.test` in root project directory and create the test database

To run test:
`make test`

## API Endpoints

| Method | Endpoint         | Description                  |
|--------|------------------|------------------------------|
| GET    | `/accounts/{account_id}`      | Get account balance                |
| POST   | `/accounts`      | Create a new account         |
| POST   | `/transactions`  | Initiate a new transaction   |

(Refer to `adapters/web/routes.go` for full routing details.)

---

## API Documentation

After running the application, access the interactive Swagger UI at:
**[http://localhost:8080/docs/index.html](http://localhost:8080/docs/index.html)**

---

## Notes

- Ensure the database is up and running before starting the application
- Ensure the test database is up and running before starting the tests
