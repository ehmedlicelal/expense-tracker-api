package main

import (
	"fmt"
	"net/http"
	"os"

	"expense-tracker/internal/database"
	"expense-tracker/internal/handler"
	"expense-tracker/internal/repository"
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

	expenseRepository := repository.NewExpenseRepository(db)
	expenseHandler := handler.NewExpenseHandler(expenseRepository)

	http.HandleFunc("/expenses", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			expenseHandler.Create(w, r)
		case http.MethodGet:
			expenseHandler.GetAll(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/expenses/summary", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		expenseHandler.GetSummary(w, r)
	})

	http.HandleFunc("/expenses/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			expenseHandler.GetByID(w, r)

		case http.MethodPatch:
			expenseHandler.Update(w, r)

		case http.MethodDelete:
			expenseHandler.Delete(w, r)

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Server is running on http://localhost:" + port)

	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println(err)
	}
}
