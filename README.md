# Personal Expense Tracker API

A simple REST API for managing personal expenses, built with Go and SQLite.

## Quick Test

### Automated Tests

Run the automated tests:

```bash
go test ./...

The tests cover expense creation, validation, retrieval, updates, deletion, not-found cases, and summary handling.

Manual API Testing

Start the server:

go run ./cmd/server

Or use Docker:

docker build -t expense-tracker .
docker run -p 8080:8080 expense-tracker
API Endpoints
POST /expenses
GET /expenses
GET /expenses/{id}
PATCH /expenses/{id}
DELETE /expenses/{id}
GET /expenses/summary
API Endpoints
POST /expenses

Creates an expense.

{
  "amount": 25.5,
  "category": "Food",
  "note": "Lunch",
  "spent_on": "2026-09-05"
}

Expected response: 201 Created

GET /expenses

Returns all expenses, newest first.

Expected response: 200 OK

GET /expenses/{id}

Returns an expense by ID.

Expected responses: 200 OK, 400 Bad Request, 404 Not Found

PATCH /expenses/{id}

Updates amount, category, and/or note.

Expected responses: 200 OK, 400 Bad Request, 404 Not Found

DELETE /expenses/{id}

Deletes an expense.

Expected responses: 200 OK, 404 Not Found

GET /expenses/summary

Returns total expenses grouped by category.

Expected response: 200 OK

Configuration

Environment variables:

PORT=8080
DATABASE_PATH=expenses.db

These values are used by default when the variables are not provided.

Technical Decisions

SQLite was chosen because it provides simple persistence without requiring a separate database server. The project uses Go's standard HTTP library, database/sql, and modernc.org/sqlite, which does not require CGO.

The code is separated into handlers, repository, database, model, and configuration layers.

Assumptions and Improvements

spent_on uses the YYYY-MM-DD format. IDs are auto-incrementing integers. Authentication is outside the assignment scope.

Possible improvements include pagination, stronger validation, structured logging, and additional repository-level tests.
