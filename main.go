package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type APIResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func writeResponse(w http.ResponseWriter, status int, message string, data any) {
	writeJSON(w, status, APIResponse{
		Message: message,
		Data:    data,
	})
}

type Expense struct {
	Id     int    `json:"id"`
	Item   string `json:"item"`
	Amount int    `json:"amount"`
}

type ExpensePatch struct {
	Item   *string `json:"item"`
	Amount *int    `json:"amount"`
}

type MemoryStore struct {
	mu     sync.Mutex
	nextID int
	items  []Expense
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{nextID: 1}
}

func (s *MemoryStore) Add(expense Expense) Expense {
	s.mu.Lock()
	defer s.mu.Unlock()

	expense.Id = s.nextID
	s.nextID++
	s.items = append(s.items, expense)

	return expense
}

func (s *MemoryStore) List() []Expense {
	s.mu.Lock()
	defer s.mu.Unlock()

	snapshot := make([]Expense, len(s.items))

	copy(snapshot, s.items)
	return snapshot
}

func (s *MemoryStore) GetByID(id int) (Expense, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, value := range s.items {

		if value.Id == id {
			return value, true
		}
	}
	return Expense{}, false
}

func getExpenseByID(store *MemoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil || id <= 0 {
			http.Error(w, "Expense ID must be a positive integer", http.StatusBadRequest)
			return
		}

		expense, found := store.GetByID(id)
		if !found {
			http.Error(w, "Expense not found", http.StatusNotFound)
			return
		}

		writeJSON(w, http.StatusOK, expense)
	}
}

func (s *MemoryStore) Update(id int, patch ExpensePatch) (Expense, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.items {
		if s.items[i].Id != id {
			continue
		}

		if patch.Item != nil {
			s.items[i].Item = *patch.Item
		}

		if patch.Amount != nil {
			s.items[i].Amount = *patch.Amount
		}

		return s.items[i], true
	}

	return Expense{}, false
}

func updateExpense(store *MemoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil || id <= 0 {
			http.Error(w, "Expense ID must be a positive integer", http.StatusBadRequest)
			return
		}

		var patch ExpensePatch
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

		// writeJSON(w, http.StatusOK, updated)

		writeResponse(w, http.StatusOK, "Expense updated successfully", updated)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	fmt.Fprintln(w, "OK")
}

/*func listExpense(w http.ResponseWriter, r *http.Request) {
	expense := []Expense{
		{Item: "Lunch", Amount: 12_500},
		{Item: "Dinner", Amount: 10_000},
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(expense); err != nil {
		log.Printf("Could not write expense response: %v", err)
	}
}

func createExpense(w http.ResponseWriter, r *http.Request) {
	var expense Expense
	if err := json.NewDecoder(r.Body).Decode(&expense); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}
	expense.Item = strings.TrimSpace(expense.Item)

	if expense.Item == "" || expense.Amount <= 0 {
		http.Error(w, "Item is required and amount must be positive", http.StatusBadRequest)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(expense); err != nil {
		log.Printf("Could not write expense response: %v", err)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("GET /expenses", listExpense)
	mux.HandleFunc("POST /expenses", createExpense)
	log.Println("Server linstening on http://127.0.0.1:8000")

	log.Fatal(http.ListenAndServe("127.0.0.1:8000", mux))
}*/

func listExpenses(store *MemoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeResponse(w, http.StatusOK, "", store.List())
		// writeJSON(w, http.StatusOK, store.List())
	}
}

func createExpense(store *MemoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var expense Expense

		if err := json.NewDecoder(r.Body).Decode(&expense); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}

		expense.Item = strings.TrimSpace(expense.Item)
		if expense.Item == "" || expense.Amount <= 0 {
			http.Error(w, "Item is required and amount must be positive", http.StatusBadRequest)
			return
		}

		expense = store.Add(expense)
		// writeJSON(w, http.StatusCreated, expense)
		writeResponse(w, http.StatusCreated, "Expense created", expense)
	}
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("Could not write JSON response: %v", err)
	}
}

func main() {
	// store := &MemoryStore{}

	store := NewMemoryStore()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("GET /expenses", listExpenses(store))
	mux.HandleFunc("POST /expenses", createExpense(store))
	mux.HandleFunc("GET /expenses/{id}", getExpenseByID(store))
	mux.HandleFunc("PATCH /expenses/{id}", updateExpense(store))

	log.Println("Server listening on http://127.0.0.1:8000")
	log.Fatal(http.ListenAndServe("127.0.0.1:8000", mux))
}
