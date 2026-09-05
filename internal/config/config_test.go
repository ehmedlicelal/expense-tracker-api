package config

import "testing"

func TestLoadDefaults(t *testing.T) {
	t.Setenv("PORT", "")
	t.Setenv("DATABASE_PATH", "")

	cfg := Load()

	if cfg.Port != "8080" {
		t.Fatalf("expected default port 8080, got %s", cfg.Port)
	}

	if cfg.DatabasePath != "expenses.db" {
		t.Fatalf("expected default database path expenses.db, got %s", cfg.DatabasePath)
	}
}

func TestLoadFromEnvironment(t *testing.T) {
	t.Setenv("PORT", "9090")
	t.Setenv("DATABASE_PATH", "test.db")

	cfg := Load()

	if cfg.Port != "9090" {
		t.Fatalf("expected port 9090, got %s", cfg.Port)
	}

	if cfg.DatabasePath != "test.db" {
		t.Fatalf("expected database path test.db, got %s", cfg.DatabasePath)
	}
}
