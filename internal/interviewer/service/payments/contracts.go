package payments

import (
	"context"
	"github.com/google/uuid"
)

type (
	CreatePaymentForSubscriptionOut struct {
		RedirectURL string
	}

	Service interface {
		CreatePaymentForSubscription(ctx context.Context, userID uuid.UUID) (*CreatePaymentForSubscriptionOut, error)
	}
)
