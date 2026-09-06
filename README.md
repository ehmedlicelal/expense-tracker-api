# Personal Expense Tracker API

A simple REST API for managing personal expenses, built with Go and SQLite.

## Quick Test

### Automated Tests

Run the complete test suite:

```bash
go test ./...
```

Tests cover model behavior, configuration, database migrations, repository CRUD operations, request validation, HTTP handlers, routing, and server integration.

### Manual Testing

Start the server:

```bash
go run ./cmd/server
```

Or use Docker:

```bash
docker build -t expense-tracker .
docker run -p 8080:8080 expense-tracker
```

## API Endpoints

* [POST /expenses](#post-expenses)
* [GET /expenses](#get-expenses)
* [GET /expenses/{id}](#get-expensesid)
* [PATCH /expenses/{id}](#patch-expensesid)
* [DELETE /expenses/{id}](#delete-expensesid)
* [GET /expenses/summary](#get-expensessummary)

### POST /expenses

Creates an expense.

```json
{
  "amount": 25.5,
  "category": "Food",
  "note": "Lunch",
  "spent_on": "2026-09-05"
}
```

Expected response: `201 Created`

### GET /expenses

Returns all expenses, newest first.

Expected response: `200 OK`

### GET /expenses/{id}

Returns an expense by ID.

Expected responses: `200 OK`, `400 Bad Request`, `404 Not Found`

### PATCH /expenses/{id}

Updates `amount`, `category`, and/or `note`.

Expected responses: `200 OK`, `400 Bad Request`, `404 Not Found`

### DELETE /expenses/{id}

Deletes an expense.

Expected responses: `200 OK`, `404 Not Found`

### GET /expenses/summary

Returns total expenses grouped by category.

Expected response: `200 OK`

## Configuration

Environment variables:

```text
PORT=8080
DATABASE_PATH=expenses.db
```

Defaults are used when these variables are not provided.

## Technical Decisions

SQLite was chosen for simple persistence without a separate database server. The project uses Go's standard HTTP library, `database/sql`, and `modernc.org/sqlite` to avoid CGO.

The code is organized into handlers, repository, database, model, and configuration layers.

## Assumptions and Improvements

`spent_on` uses the `YYYY-MM-DD` format. IDs are auto-incrementing integers. Authentication is outside the assignment scope.

Possible improvements include pagination, stronger validation, structured logging, and improved database handling for larger workloads.
