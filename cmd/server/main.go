package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/expenses", expensesHandler)

	fmt.Println("Server is running on http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)
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

	fmt.Fprintln(w, "Expenses endpoint")
}
