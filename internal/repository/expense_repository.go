package repository

import (
	"database/sql"
	"expense-tracker/internal/expense"
)

type ExpenseRepository struct {
	db *sql.DB
}

func NewExpenseRepository(db *sql.DB) *ExpenseRepository {
	return &ExpenseRepository{
		db: db,
	}
}

func (r *ExpenseRepository) Create(e *expense.Expense) error {
	result, err := r.db.Exec(`
		INSERT INTO expenses (amount, category, note, spent_on)
		VALUES (?, ?, ?, ?)
	`, e.Amount, e.Category, e.Note, e.SpentOn)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	e.ID = int(id)

	return r.db.QueryRow(`
		SELECT created_at
		FROM expenses
		WHERE id = ?
	`, e.ID).Scan(&e.CreatedAt)
}

func (r *ExpenseRepository) GetAll() ([]expense.Expense, error) {
	rows, err := r.db.Query(`
		SELECT id, amount, category, note, spent_on, created_at
		FROM expenses
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	expenses := make([]expense.Expense, 0)

	for rows.Next() {
		var e expense.Expense

		err := rows.Scan(
			&e.ID,
			&e.Amount,
			&e.Category,
			&e.Note,
			&e.SpentOn,
			&e.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		expenses = append(expenses, e)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return expenses, nil
}

func (r *ExpenseRepository) GetByID(id int) (*expense.Expense, error) {
	var e expense.Expense

	err := r.db.QueryRow(`
		SELECT id, amount, category, note, spent_on, created_at
		FROM expenses
		WHERE id = ?
	`, id).Scan(
		&e.ID,
		&e.Amount,
		&e.Category,
		&e.Note,
		&e.SpentOn,
		&e.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &e, nil
}

func (r *ExpenseRepository) Update(e *expense.Expense) error {
	_, err := r.db.Exec(`
		UPDATE expenses
		SET amount = ?, category = ?, note = ?
		WHERE id = ?
	`, e.Amount, e.Category, e.Note, e.ID)

	return err
}

func (r *ExpenseRepository) Delete(id int) error {
	_, err := r.db.Exec(`
		DELETE FROM expenses
		WHERE id = ?
	`, id)

	return err
}

func (r *ExpenseRepository) GetSummary() (map[string]float64, error) {
	rows, err := r.db.Query(`
		SELECT category, SUM(amount)
		FROM expenses
		GROUP BY category
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	summary := make(map[string]float64)

	for rows.Next() {
		var category string
		var total float64

		err := rows.Scan(&category, &total)
		if err != nil {
			return nil, err
		}

		summary[category] = total
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return summary, nil
}
