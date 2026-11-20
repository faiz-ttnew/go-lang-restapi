package api

import (
	"net/http"

	"go-lang-restapi/expense"
)

func NewRouter(store *expense.MemoryStore) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("GET /expenses", listExpenses(store))
	mux.HandleFunc("POST /expenses", createExpense(store))
	mux.HandleFunc("GET /expenses/{id}", getExpenseByID(store))
	mux.HandleFunc("PATCH /expenses/{id}", updateExpense(store))

	return mux
}
