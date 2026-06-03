package domain

type Order struct {
	ID        string
	ProductID string
	Quantity  int32
	Status    string // "PENDING", "RESERVED", "FAILED"
}
