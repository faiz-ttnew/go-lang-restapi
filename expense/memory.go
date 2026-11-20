package expense

import "sync"

type MemoryStore struct {
	mu     sync.Mutex
	nextID int
	items  []Expense
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{nextID: 1}
}

func (s *MemoryStore) Add(item Expense) Expense {
	s.mu.Lock()
	defer s.mu.Unlock()

	item.Id = s.nextID
	s.nextID++
	s.items = append(s.items, item)

	return item
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

	for _, item := range s.items {
		if item.Id == id {
			return item, true
		}
	}

	return Expense{}, false
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
