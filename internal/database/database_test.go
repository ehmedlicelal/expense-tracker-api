package database

import "testing"

func TestOpen(t *testing.T) {
	db, err := Open(":memory:")
	if err != nil {
		t.Fatalf("expected database to open, got error: %v", err)
	}

	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Fatalf("expected database ping to succeed, got error: %v", err)
	}
}

func TestRunMigrations(t *testing.T) {
	db, err := Open(":memory:")
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	defer db.Close()

	if err := RunMigrations(
		db,
		"../../migrations/001_create_expenses.sql",
	); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	var tableName string

	err = db.QueryRow(`
		SELECT name
		FROM sqlite_master
		WHERE type = 'table' AND name = 'expenses'
	`).Scan(&tableName)

	if err != nil {
		t.Fatalf("expected expenses table to exist: %v", err)
	}

	if tableName != "expenses" {
		t.Fatalf("expected expenses table, got %s", tableName)
	}
}

func TestRunMigrationsIsIdempotent(t *testing.T) {
	db, err := Open(":memory:")
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	defer db.Close()

	if err := RunMigrations(
		db,
		"../../migrations/001_create_expenses.sql",
	); err != nil {
		t.Fatalf("first migration failed: %v", err)
	}

	if err := RunMigrations(
		db,
		"../../migrations/001_create_expenses.sql",
	); err != nil {
		t.Fatalf("second migration failed: %v", err)
	}

	var count int

	err = db.QueryRow(`
		SELECT COUNT(*)
		FROM sqlite_master
		WHERE type = 'table' AND name = 'expenses'
	`).Scan(&count)

	if err != nil {
		t.Fatalf("failed to check expenses table: %v", err)
	}

	if count != 1 {
		t.Fatalf("expected exactly one expenses table, got %d", count)
	}
}
