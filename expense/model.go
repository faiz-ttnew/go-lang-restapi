package expense

type Expense struct {
	Id     int    `json:"id"`
	Item   string `json:"item"`
	Amount int    `json:"amount"`
}

type ExpensePatch struct {
	Item   *string `json:"item"`
	Amount *int    `json:"amount"`
}
