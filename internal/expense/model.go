package expense

import "time"

type Expense struct {
	ID        int       `json:"id"`
	Amount    float64   `json:"amount"`
	Category  string    `json:"category"`
	Note      string    `json:"note,omitempty"`
	SpentOn   time.Time `json:"spent_on"`
	CreatedAt time.Time `json:"created_at"`
}
