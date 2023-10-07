package payments

import (
	"context"
	"github.com/google/uuid"
	"job-interviewer/pkg/payments/model"
	"job-interviewer/pkg/transactional"
)

type (
	CreatePaymentIn struct {
		IDK         uuid.UUID
		UserID      uuid.UUID
		Type        model.Type
		Description string
		Amount      model.Penny
	}

	HandlePaymentCompletedIn struct {
		ExternalID model.ExternalID
		Status     model.Status
	}

	PaymentResult struct {
		PaymentType model.Type
		Paid        bool
	}

	Service interface {
		CreatePayment(ctx context.Context, in *CreatePaymentIn) (*model.Payment, error)
		CheckPendingPayment(ctx context.Context, userID uuid.UUID) error
		CheckPendingPaymentWithResult(ctx context.Context, userID uuid.UUID) (*PaymentResult, error)

		RegisterPaymentPaidHandler(paymentType model.Type, handler PaymentCompletedHandler) error
		RegisterPaymentCanceledHandler(paymentType model.Type, handler PaymentCompletedHandler) error
	}

	PaymentCompletedHandler interface {
		Handle(ctx context.Context, userID uuid.UUID) error
	}

	repository interface {
		CreatePayment(ctx context.Context, tx transactional.Tx, in *model.Payment) error
		UpdatePayment(ctx context.Context, tx transactional.Tx, in *model.Payment) error
		GetPaymentByExternalID(ctx context.Context, ID model.ExternalID) (*model.Payment, error)
		GetActivePaymentByUserID(ctx context.Context, userID uuid.UUID) (*model.Payment, error)
	}
)
