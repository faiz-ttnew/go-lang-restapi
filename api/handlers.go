package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"go-lang-restapi/expense"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintln(w, "OK")
}

func listExpenses(store *expense.MemoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeResponse(w, http.StatusOK, "", store.List())
	}
}

func createExpense(store *expense.MemoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var item expense.Expense

		if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}

		item.Item = strings.TrimSpace(item.Item)
		if item.Item == "" || item.Amount <= 0 {
			http.Error(w, "Item is required and amount must be positive", http.StatusBadRequest)
			return
		}

		created := store.Add(item)
		writeResponse(w, http.StatusCreated, "Expense created", created)
	}
}

func getExpenseByID(store *expense.MemoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil || id <= 0 {
			http.Error(w, "Expense ID must be a positive integer", http.StatusBadRequest)
			return
		}

		item, found := store.GetByID(id)
		if !found {
			http.Error(w, "Expense not found", http.StatusNotFound)
			return
		}

		writeJSON(w, http.StatusOK, item)
	}
}

func updateExpense(store *expense.MemoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil || id <= 0 {
			http.Error(w, "Expense ID must be a positive integer", http.StatusBadRequest)
			return
		}

		var patch expense.ExpensePatch
		if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}

		if patch.Item == nil && patch.Amount == nil {
			http.Error(w, "Provide item or amount", http.StatusBadRequest)
			return
		}

		if patch.Item != nil {
			item := strings.TrimSpace(*patch.Item)
			if item == "" {
				http.Error(w, "Item cannot be blank", http.StatusBadRequest)
				return
			}
			patch.Item = &item
		}

		if patch.Amount != nil && *patch.Amount <= 0 {
			http.Error(w, "Amount must be positive", http.StatusBadRequest)
			return
		}

		updated, found := store.Update(id, patch)
		if !found {
			http.Error(w, "Expense not found", http.StatusNotFound)
			return
		}

		writeResponse(w, http.StatusOK, "Expense updated successfully", updated)
	}
}
