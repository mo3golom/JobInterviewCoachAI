package payments

import (
	"context"
	"errors"
	"github.com/google/uuid"
)

var (
	ErrHandlerAlreadyRegistered = errors.New("handler is already registered for this type")
	ErrPaymentWrongStatus       = errors.New("payment has a wrong status for this action")
	ErrPaymentNotFound          = errors.New("payment not found")
)

type (
	CreatePaymentIn struct {
		IDK         uuid.UUID
		UserID      uuid.UUID
		Type        Type
		Description string
		Amount      Penny
	}

	HandlePaymentCompletedIn struct {
		ExternalID ExternalID
		Status     Status
	}

	Service interface {
		CreatePayment(ctx context.Context, in *CreatePaymentIn) (*Payment, error)
		CheckPendingPayment(ctx context.Context, userID uuid.UUID) error

		RegisterPaymentPaidHandler(paymentType Type, handler PaymentCompletedHandler) error
		RegisterPaymentCanceledHandler(paymentType Type, handler PaymentCompletedHandler) error
	}

	PaymentCompletedHandler interface {
		Handle(ctx context.Context, userID uuid.UUID) error
	}

	GatewayCreatePaymentIn struct {
		IDK         uuid.UUID
		Description string
		Amount      int64
	}

	GatewayCreatePaymentOut struct {
		ExternalID  ExternalID
		RedirectURl string
	}

	Gateway interface {
		CreatePayment(ctx context.Context, in *GatewayCreatePaymentIn) (*GatewayCreatePaymentOut, error)
		GetPaymentStatus(ctx context.Context, ID ExternalID) (*Status, error)
	}

	repository interface {
		CreatePayment(ctx context.Context, in *NewPayment) error
		UpdatePayment(ctx context.Context, in *Payment) error
		GetPaymentByExternalID(ctx context.Context, ID ExternalID) (*Payment, error)
		GetPendingPaymentByUserID(ctx context.Context, userID uuid.UUID) (*Payment, error)
	}
)
