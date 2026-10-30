package models

type Category struct {
	ID   int64
	Name string
}

type Expense struct {
	ID          int64
	Amount      float64
	Description string
	Date        string
	CategoryID  int64
	Category    string
}
