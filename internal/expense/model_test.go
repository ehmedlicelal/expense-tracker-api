package expense

import (
	"encoding/json"
	"testing"
	"time"
)

func TestExpenseJSON(t *testing.T) {
	spentOn := time.Date(2026, 9, 5, 0, 0, 0, 0, time.UTC)
	createdAt := time.Date(2026, 9, 5, 12, 30, 0, 0, time.UTC)

	e := Expense{
		ID:        1,
		Amount:    25.50,
		Category:  "Food",
		Note:      "Lunch",
		SpentOn:   spentOn,
		CreatedAt: createdAt,
	}

	data, err := json.Marshal(e)
	if err != nil {
		t.Fatalf("failed to marshal expense: %v", err)
	}

	var result map[string]interface{}

	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	if result["id"] != float64(1) {
		t.Fatalf("expected id 1, got %v", result["id"])
	}

	if result["amount"] != 25.50 {
		t.Fatalf("expected amount 25.50, got %v", result["amount"])
	}

	if result["category"] != "Food" {
		t.Fatalf("expected category Food, got %v", result["category"])
	}

	if result["note"] != "Lunch" {
		t.Fatalf("expected note Lunch, got %v", result["note"])
	}
}

func TestExpenseJSONOmitsEmptyNote(t *testing.T) {
	e := Expense{
		ID:       1,
		Amount:   10,
		Category: "Food",
	}

	data, err := json.Marshal(e)
	if err != nil {
		t.Fatalf("failed to marshal expense: %v", err)
	}

	var result map[string]interface{}

	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	if _, exists := result["note"]; exists {
		t.Fatal("expected empty note to be omitted from JSON")
	}
}
