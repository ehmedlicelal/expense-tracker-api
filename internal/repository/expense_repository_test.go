package repository

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"expense-tracker/internal/database"
	"expense-tracker/internal/expense"

	_ "modernc.org/sqlite"
)

func setupTestRepository(t *testing.T) *ExpenseRepository {
	t.Helper()

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	db.SetMaxOpenConns(1)

	err = database.RunMigrations(
		db,
		"../../migrations/001_create_expenses.sql",
	)
	if err != nil {
		t.Fatal(err)
	}

	return NewExpenseRepository(db)
}

func createTestExpense(
	t *testing.T,
	repo *ExpenseRepository,
	amount float64,
	category string,
) *expense.Expense {
	t.Helper()

	e := &expense.Expense{
		Amount:   amount,
		Category: category,
		Note:     "Test expense",
		SpentOn:  time.Date(2026, 9, 5, 0, 0, 0, 0, time.UTC),
	}

	if err := repo.Create(e); err != nil {
		t.Fatal(err)
	}

	return e
}

func TestCreate(t *testing.T) {
	repo := setupTestRepository(t)

	e := createTestExpense(t, repo, 25.50, "Food")

	if e.ID == 0 {
		t.Fatal("expected ID to be assigned")
	}

	if e.CreatedAt.IsZero() {
		t.Fatal("expected CreatedAt to be assigned")
	}
}

func TestGetByID(t *testing.T) {
	repo := setupTestRepository(t)

	created := createTestExpense(t, repo, 25.50, "Food")

	found, err := repo.GetByID(created.ID)
	if err != nil {
		t.Fatal(err)
	}

	if found.ID != created.ID {
		t.Fatalf("expected ID %d, got %d", created.ID, found.ID)
	}

	if found.Amount != 25.50 {
		t.Fatalf("expected amount 25.50, got %v", found.Amount)
	}

	if found.Category != "Food" {
		t.Fatalf("expected category Food, got %s", found.Category)
	}

	if found.Note != "Test expense" {
		t.Fatalf("expected note %q, got %q", "Test expense", found.Note)
	}
}

func TestGetByIDNotFound(t *testing.T) {
	repo := setupTestRepository(t)

	_, err := repo.GetByID(999)

	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("expected sql.ErrNoRows, got %v", err)
	}
}

func TestGetAll(t *testing.T) {
	repo := setupTestRepository(t)

	first := createTestExpense(t, repo, 10, "Food")
	second := createTestExpense(t, repo, 20, "Transport")

	expenses, err := repo.GetAll()
	if err != nil {
		t.Fatal(err)
	}

	if len(expenses) != 2 {
		t.Fatalf("expected 2 expenses, got %d", len(expenses))
	}

	if expenses[0].ID != second.ID {
		t.Fatalf("expected newest expense first, got ID %d", expenses[0].ID)
	}

	if expenses[1].ID != first.ID {
		t.Fatalf("expected oldest expense second, got ID %d", expenses[1].ID)
	}
}

func TestUpdate(t *testing.T) {
	repo := setupTestRepository(t)

	created := createTestExpense(t, repo, 15, "Food")

	updatedExpense := &expense.Expense{
		ID:       created.ID,
		Amount:   30,
		Category: "Transport",
		Note:     "Updated note",
	}

	err := repo.Update(updatedExpense)
	if err != nil {
		t.Fatal(err)
	}

	updated, err := repo.GetByID(created.ID)
	if err != nil {
		t.Fatal(err)
	}

	if updated.Amount != 30 {
		t.Fatalf("expected amount 30, got %v", updated.Amount)
	}

	if updated.Category != "Transport" {
		t.Fatalf("expected category Transport, got %s", updated.Category)
	}

	if updated.Note != "Updated note" {
		t.Fatalf("expected note %q, got %q", "Updated note", updated.Note)
	}
}

func TestDelete(t *testing.T) {
	repo := setupTestRepository(t)

	created := createTestExpense(t, repo, 15, "Food")

	if err := repo.Delete(created.ID); err != nil {
		t.Fatal(err)
	}

	_, err := repo.GetByID(created.ID)

	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("expected sql.ErrNoRows after delete, got %v", err)
	}
}

func TestGetSummary(t *testing.T) {
	repo := setupTestRepository(t)

	createTestExpense(t, repo, 20, "Food")
	createTestExpense(t, repo, 30, "Food")
	createTestExpense(t, repo, 15, "Transport")
	createTestExpense(t, repo, 25, "Transport")

	summary, err := repo.GetSummary()
	if err != nil {
		t.Fatal(err)
	}

	if summary["Food"] != 50 {
		t.Fatalf("expected Food total 50, got %v", summary["Food"])
	}

	if summary["Transport"] != 40 {
		t.Fatalf("expected Transport total 40, got %v", summary["Transport"])
	}
}
