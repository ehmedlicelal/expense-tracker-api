package handler

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"expense-tracker/internal/database"
	"expense-tracker/internal/repository"
)

func setupTestHandler(t *testing.T) *ExpenseHandler {
	t.Helper()

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	err = database.RunMigrations(db, "../../migrations/001_create_expenses.sql")
	if err != nil {
		t.Fatal(err)
	}

	repo := repository.NewExpenseRepository(db)

	return NewExpenseHandler(repo)
}

func TestCreateExpense(t *testing.T) {
	handler := setupTestHandler(t)

	body := `{
		"amount": 25.5,
		"category": "Food",
		"note": "Lunch",
		"spent_on": "2026-09-05"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/expenses",
		strings.NewReader(body),
	)

	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}
}

func TestCreateExpenseInvalidAmount(t *testing.T) {
	handler := setupTestHandler(t)

	body := `{
		"amount": -10,
		"category": "Food",
		"spent_on": "2026-09-05"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/expenses",
		strings.NewReader(body),
	)

	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestCreateExpenseEmptyCategory(t *testing.T) {
	handler := setupTestHandler(t)

	body := `{
		"amount": 25,
		"category": "",
		"spent_on": "2026-09-05"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/expenses",
		strings.NewReader(body),
	)

	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestCreateExpenseInvalidDate(t *testing.T) {
	handler := setupTestHandler(t)

	body := `{
		"amount": 25,
		"category": "Food",
		"spent_on": "05-09-2026"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/expenses",
		strings.NewReader(body),
	)

	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestGetExpenseByID(t *testing.T) {
	handler := setupTestHandler(t)

	// First create an expense
	body := `{
		"amount": 50,
		"category": "Transport",
		"note": "Taxi",
		"spent_on": "2026-09-05"
	}`

	createReq := httptest.NewRequest(
		http.MethodPost,
		"/expenses",
		strings.NewReader(body),
	)

	createRec := httptest.NewRecorder()

	handler.Create(createRec, createReq)

	if createRec.Code != http.StatusCreated {
		t.Fatalf("expected create status %d, got %d", http.StatusCreated, createRec.Code)
	}

	// Get the created expense
	getReq := httptest.NewRequest(
		http.MethodGet,
		"/expenses/1",
		nil,
	)

	getRec := httptest.NewRecorder()

	handler.GetByID(getRec, getReq)

	if getRec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, getRec.Code)
	}
}

func TestGetExpenseByIDNotFound(t *testing.T) {
	handler := setupTestHandler(t)

	req := httptest.NewRequest(
		http.MethodGet,
		"/expenses/999",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.GetByID(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestGetAllExpenses(t *testing.T) {
	handler := setupTestHandler(t)

	body := `{
		"amount": 30,
		"category": "Food",
		"spent_on": "2026-09-05"
	}`

	createReq := httptest.NewRequest(
		http.MethodPost,
		"/expenses",
		strings.NewReader(body),
	)

	createRec := httptest.NewRecorder()

	handler.Create(createRec, createReq)

	if createRec.Code != http.StatusCreated {
		t.Fatalf("expected create status %d, got %d", http.StatusCreated, createRec.Code)
	}

	// Get all expenses
	getReq := httptest.NewRequest(
		http.MethodGet,
		"/expenses",
		nil,
	)

	getRec := httptest.NewRecorder()

	handler.GetAll(getRec, getReq)

	if getRec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, getRec.Code)
	}
}

func TestUpdateExpense(t *testing.T) {
	handler := setupTestHandler(t)

	// Create an expense first
	createBody := `{
		"amount": 50,
		"category": "Food",
		"note": "Lunch",
		"spent_on": "2026-09-05"
	}`

	createReq := httptest.NewRequest(
		http.MethodPost,
		"/expenses",
		strings.NewReader(createBody),
	)

	createRec := httptest.NewRecorder()

	handler.Create(createRec, createReq)

	if createRec.Code != http.StatusCreated {
		t.Fatalf("expected create status %d, got %d", http.StatusCreated, createRec.Code)
	}

	// Update the expense
	updateBody := `{
		"amount": 75,
		"category": "Transport"
	}`

	updateReq := httptest.NewRequest(
		http.MethodPatch,
		"/expenses/1",
		strings.NewReader(updateBody),
	)

	updateRec := httptest.NewRecorder()

	handler.Update(updateRec, updateReq)

	if updateRec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, updateRec.Code)
	}
}

func TestUpdateExpenseInvalidAmount(t *testing.T) {
	handler := setupTestHandler(t)

	createBody := `{
		"amount": 50,
		"category": "Food",
		"spent_on": "2026-09-05"
	}`

	createReq := httptest.NewRequest(
		http.MethodPost,
		"/expenses",
		strings.NewReader(createBody),
	)

	createRec := httptest.NewRecorder()

	handler.Create(createRec, createReq)

	updateBody := `{
		"amount": -20
	}`

	updateReq := httptest.NewRequest(
		http.MethodPatch,
		"/expenses/1",
		strings.NewReader(updateBody),
	)

	updateRec := httptest.NewRecorder()

	handler.Update(updateRec, updateReq)

	if updateRec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, updateRec.Code)
	}
}

func TestDeleteExpense(t *testing.T) {
	handler := setupTestHandler(t)

	createBody := `{
		"amount": 50,
		"category": "Food",
		"spent_on": "2026-09-05"
	}`

	createReq := httptest.NewRequest(
		http.MethodPost,
		"/expenses",
		strings.NewReader(createBody),
	)

	createRec := httptest.NewRecorder()

	handler.Create(createRec, createReq)

	if createRec.Code != http.StatusCreated {
		t.Fatalf("expected create status %d, got %d", http.StatusCreated, createRec.Code)
	}

	deleteReq := httptest.NewRequest(
		http.MethodDelete,
		"/expenses/1",
		nil,
	)

	deleteRec := httptest.NewRecorder()

	handler.Delete(deleteRec, deleteReq)

	if deleteRec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, deleteRec.Code)
	}

	// Make sure the expense was actually deleted
	getReq := httptest.NewRequest(
		http.MethodGet,
		"/expenses/1",
		nil,
	)

	getRec := httptest.NewRecorder()

	handler.GetByID(getRec, getReq)

	if getRec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d after delete, got %d", http.StatusNotFound, getRec.Code)
	}
}

func TestGetExpenseSummary(t *testing.T) {
	handler := setupTestHandler(t)

	expenses := []string{
		`{
			"amount": 25,
			"category": "Food",
			"spent_on": "2026-09-05"
		}`,
		`{
			"amount": 15,
			"category": "Food",
			"spent_on": "2026-09-05"
		}`,
		`{
			"amount": 40,
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

		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("expected create status %d, got %d", http.StatusCreated, rec.Code)
		}
	}

	req := httptest.NewRequest(
		http.MethodGet,
		"/expenses/summary",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.GetSummary(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	responseBody := rec.Body.String()

	if !strings.Contains(responseBody, `"Food":40`) {
		t.Fatalf("expected Food total to be 40, got %s", responseBody)
	}

	if !strings.Contains(responseBody, `"Transport":40`) {
		t.Fatalf("expected Transport total to be 40, got %s", responseBody)
	}
}
