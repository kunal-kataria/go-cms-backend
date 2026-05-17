# Content Management System Backend

Content management systems (CMS) help manage and deliver content across websites and applications. This project builds a CMS backend API using Go, Gin, and PostgreSQL, covering CRUD operations, migrations, testing, and environment-based configuration.

## Project Deliverables & Outcomes

- Developed scalable REST APIs for pages, posts, and media using Gin and Go.
- Integrated PostgreSQL with GORM for database operations and ORM workflows.
- Managed database migrations and schema updates.
- Implemented unit and integration testing for reliability.
- Configured secure environment-based setups for development and production.

## Prerequisites

- Go 1.20+ (or compatible)
- PostgreSQL 13+
- golang-migrate

## Installation steps

1. **Clone the repository**

    ```bash
    git clone https://github.com/kunal-kataria/go-cms-backend.git
    cd go-cms-backend
    ```

2. **Install Go dependencies**

    ```bash
    go mod download
    ```

3. **Set Up Environment Variables**

    Create a `.env` file. In the project root directory, create a `.env` file to store your development environment variables:

    ```bash
    cp .env.example .env
    ```

    Open the `.env` file and replace the placeholder values with your actual database credentials:
    - `DB_USER`
    - `DB_PASSWORD`
    - `DB_NAME`
    - `DB_HOST`
    - `DB_PORT`
    - `ENV` (optional: `development` or `production`, default is `development`)

    For integration tests, set a separate test database:
    - `TEST_DB_USER`
    - `TEST_DB_PASSWORD`
    - `TEST_DB_NAME`
    - `TEST_DB_HOST`
    - `TEST_DB_PORT`

## Migrations

```bash
migrate -database "postgres://username:password@localhost:port/dbname?sslmode=disable" -path ./migrations up
```

Rollback migrations:

```bash
migrate -database "postgres://username:password@localhost:port/dbname?sslmode=disable" -path ./migrations down
```

## Run

```bash
go run main.go
```

The API will start on `:8080`.

## Routes

Base path: `/api/v1`

- Pages: `GET /pages`, `GET /pages/:id`, `POST /pages`, `PUT /pages/:id`, `DELETE /pages/:id`
- Posts: `GET /posts`, `GET /posts/:id`, `POST /posts`, `PUT /posts/:id`, `DELETE /posts/:id`
- Media: `GET /media`, `GET /media/:id`, `POST /media`, `DELETE /media/:id`

## Tests

```bash
# Run all unit tests with verbose output
go test ./controllers -v

# Run specific controller tests
go test ./controllers -run TestGetMedia -v
go test ./controllers -run TestCreatePost -v
go test ./controllers -run TestUpdatePage -v

# Run tests with coverage
go test ./controllers -cover

# Run all integration tests
go test ./test/integration
```
