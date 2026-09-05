package database

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

func Open(databasePath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", databasePath)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func RunMigrations(db *sql.DB, migrationPath string) error {
	migration, err := os.ReadFile(migrationPath)
	if err != nil {
		return err
	}

	_, err = db.Exec(string(migration))
	if err != nil {
		return err
	}

	return nil
}
