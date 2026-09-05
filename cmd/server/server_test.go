package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"expense-tracker/internal/database"
	"expense-tracker/internal/handler"
	"expense-tracker/internal/repository"
)

func setupTestServer(t *testing.T) http.Handler {
	t.Helper()

	db, err := database.Open(":memory:")
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	db.SetMaxOpenConns(1)

	if err := database.RunMigrations(
		db,
		"../../migrations/001_create_expenses.sql",
	); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	expenseRepository := repository.NewExpenseRepository(db)
	expenseHandler := handler.NewExpenseHandler(expenseRepository)

	return setupRoutes(expenseHandler)
}

func TestServerCreateAndGetExpense(t *testing.T) {
	server := setupTestServer(t)

	body := `{
		"amount": 25.50,
		"category": "Food",
		"note": "Lunch",
		"spent_on": "2026-09-05"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/expenses",
		strings.NewReader(body),
	)

	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	server.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusCreated {
		t.Fatalf(
			"expected status 201, got %d",
			recorder.Code,
		)
	}

	getRequest := httptest.NewRequest(
		http.MethodGet,
		"/expenses/1",
		nil,
	)

	getRecorder := httptest.NewRecorder()

	server.ServeHTTP(getRecorder, getRequest)

	if getRecorder.Code != http.StatusOK {
		t.Fatalf(
			"expected status 200, got %d",
			getRecorder.Code,
		)
	}

	if !strings.Contains(getRecorder.Body.String(), `"Food"`) {
		t.Fatalf("expected response to contain Food")
	}
}

func TestServerSummary(t *testing.T) {
	server := setupTestServer(t)

	expenses := []string{
		`{
			"amount": 25,
			"category": "Food",
			"spent_on": "2026-09-05"
		}`,
		`{
			"amount": 15,
			"category": "Transport",
			"spent_on": "2026-09-05"
		}`,
	}

	for _, body := range expenses {
		req := httptest.NewRequest(
			http.MethodPost,
			"/expenses",
			strings.NewReader(body),
		)

		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		server.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusCreated {
			t.Fatalf(
				"expected status 201, got %d",
				recorder.Code,
			)
		}
	}

	req := httptest.NewRequest(
		http.MethodGet,
		"/expenses/summary",
		nil,
	)

	recorder := httptest.NewRecorder()

	server.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf(
			"expected status 200, got %d",
			recorder.Code,
		)
	}

	response := recorder.Body.String()

	if !strings.Contains(response, `"Food":25`) {
		t.Fatalf("expected Food total to be 25, got %s", response)
	}

	if !strings.Contains(response, `"Transport":15`) {
		t.Fatalf("expected Transport total to be 15, got %s", response)
	}
}

func TestServerInvalidMethod(t *testing.T) {
	server := setupTestServer(t)

	req := httptest.NewRequest(
		http.MethodPut,
		"/expenses",
		nil,
	)

	recorder := httptest.NewRecorder()

	server.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusMethodNotAllowed {
		t.Fatalf(
			"expected status 405, got %d",
			recorder.Code,
		)
	}
}
