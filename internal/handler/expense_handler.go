package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"expense-tracker/internal/expense"
	"expense-tracker/internal/repository"
)

type createExpenseRequest struct {
	Amount   float64 `json:"amount"`
	Category string  `json:"category"`
	Note     string  `json:"note"`
	SpentOn  string  `json:"spent_on"`
}

type updateExpenseRequest struct {
	Amount   *float64 `json:"amount"`
	Category *string  `json:"category"`
	Note     *string  `json:"note"`
}

type ExpenseHandler struct {
	repository *repository.ExpenseRepository
}

func NewExpenseHandler(repo *repository.ExpenseRepository) *ExpenseHandler {
	return &ExpenseHandler{
		repository: repo,
	}
}

func (h *ExpenseHandler) Create(w http.ResponseWriter, r *http.Request) {
	var request createExpenseRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if request.Amount <= 0 {
		http.Error(w, "amount must be greater than 0", http.StatusBadRequest)
		return
	}

	if request.Category == "" {
		http.Error(w, "category is required", http.StatusBadRequest)
		return
	}

	spentOn, err := time.Parse("2006-01-02", request.SpentOn)
	if err != nil {
		http.Error(w, "spent_on must be in YYYY-MM-DD format", http.StatusBadRequest)
		return
	}

	e := expense.Expense{
		Amount:   request.Amount,
		Category: request.Category,
		Note:     request.Note,
		SpentOn:  spentOn,
	}

	if err := h.repository.Create(&e); err != nil {
		http.Error(w, "failed to create expense", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(e)
}

func (h *ExpenseHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	expenses, err := h.repository.GetAll()
	if err != nil {
		http.Error(w, "failed to get expenses", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(expenses)
}

func (h *ExpenseHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/expenses/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid expense id", http.StatusBadRequest)
		return
	}

	e, err := h.repository.GetByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "expense not found", http.StatusNotFound)
			return
		}

		http.Error(w, "failed to get expense", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(e)
}

func (h *ExpenseHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/expenses/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid expense id", http.StatusBadRequest)
		return
	}

	var request updateExpenseRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if request.Amount == nil && request.Category == nil && request.Note == nil {
		http.Error(w, "at least one field is required", http.StatusBadRequest)
		return
	}

	e, err := h.repository.GetByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "expense not found", http.StatusNotFound)
			return
		}

		http.Error(w, "failed to get expense", http.StatusInternalServerError)
		return
	}

	if request.Amount != nil {
		if *request.Amount <= 0 {
			http.Error(w, "amount must be greater than 0", http.StatusBadRequest)
			return
		}

		e.Amount = *request.Amount
	}

	if request.Category != nil {
		if *request.Category == "" {
			http.Error(w, "category cannot be empty", http.StatusBadRequest)
			return
		}

		e.Category = *request.Category
	}

	if request.Note != nil {
		e.Note = *request.Note
	}

	if err := h.repository.Update(e); err != nil {
		http.Error(w, "failed to update expense", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(e)
}

func (h *ExpenseHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/expenses/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid expense id", http.StatusBadRequest)
		return
	}

	_, err = h.repository.GetByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "expense not found", http.StatusNotFound)
			return
		}

		http.Error(w, "failed to get expense", http.StatusInternalServerError)
		return
	}

	if err := h.repository.Delete(id); err != nil {
		http.Error(w, "failed to delete expense", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *ExpenseHandler) GetSummary(w http.ResponseWriter, r *http.Request) {
	summary, err := h.repository.GetSummary()
	if err != nil {
		http.Error(w, "failed to get expense summary", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(summary)
}
