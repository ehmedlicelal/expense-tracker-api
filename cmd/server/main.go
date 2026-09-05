package main

import (
	"fmt"
	"net/http"
	"os"

	"expense-tracker/internal/database"
)

func main() {
	databasePath := os.Getenv("DATABASE_PATH")

	if databasePath == "" {
		databasePath = "expenses.db"
	}

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	db, err := database.Open(databasePath)
	if err != nil {
		fmt.Println("failed to open database:", err)
		return
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		fmt.Println("failed to ping database:", err)
		return
	}

	fmt.Println("Database connected successfully")

	if err := database.RunMigrations(db, "migrations/001_create_expenses.sql"); err != nil {
		fmt.Println("failed to run migrations:", err)
		return
	}

	fmt.Println("Migrations completed successfully")

	http.HandleFunc("/expenses", expensesHandler)

	fmt.Println("Server is running on http://localhost:" + port)

	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println(err)
	}
}

func expensesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	fmt.Fprintln(w, `{"message":"Expenses endpoint"}`)
}
