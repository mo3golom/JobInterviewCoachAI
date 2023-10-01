package payments

import "github.com/google/uuid"

const (
	StatusNew      Status = "new"
	StatusPending  Status = "pending"
	StatusPaid     Status = "paid"
	StatusCanceled Status = "canceled"
)

type (
	Type   string
	Status string
	Penny  int64

	ExternalID string

	NewPayment struct {
		ID          uuid.UUID
		UserID      uuid.UUID
		Amount      Penny
		Type        Type
		Description string
	}

	Payment struct {
		ID          uuid.UUID
		ExternalID  ExternalID
		UserID      uuid.UUID
		Amount      Penny
		Type        Type
		Description string
		RedirectURL string
		Status      Status
	}
)

// Normalize raises order, for example 100 penny = 1 dollar.
// I don't know which the right way to called this.
func (p Penny) Normalize() int64 {
	return int64(p / 100)
}
