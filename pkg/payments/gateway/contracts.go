package gateway

import (
	"context"
	"github.com/google/uuid"
	"job-interviewer/pkg/payments/model"
)

type (
	CreatePaymentIn struct {
		IDK         uuid.UUID
		Description string
		Amount      int64
	}

	CreatePaymentOut struct {
		ExternalID  model.ExternalID
		RedirectURl string
	}

	Gateway interface {
		CreatePayment(ctx context.Context, in *CreatePaymentIn) (*CreatePaymentOut, error)
		GetPaymentStatus(ctx context.Context, ID model.ExternalID) (*model.Status, error)
	}
)
