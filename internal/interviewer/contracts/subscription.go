package contracts

import (
	"context"
	"github.com/google/uuid"
)

type (
	CreatePaymentOut struct {
		RedirectURL string
	}

	SubscriptionUseCase interface {
		CreatePayment(ctx context.Context, userID uuid.UUID) (*CreatePaymentOut, error)
	}
)
