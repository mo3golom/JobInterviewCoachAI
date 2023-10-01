package subscription

import (
	"context"
	"github.com/google/uuid"
	"job-interviewer/internal/interviewer/contracts"
	"job-interviewer/internal/interviewer/service/payments"
)

type UseCase struct {
	payments payments.Service
}

func (u *UseCase) CreatePayment(ctx context.Context, userID uuid.UUID) (*contracts.CreatePaymentOut, error) {
	result, err := u.payments.CreatePaymentForSubscription(
		ctx,
		userID,
	)
	if err != nil {
		return nil, err
	}

	return &contracts.CreatePaymentOut{
		RedirectURL: result.RedirectURL,
	}, nil
}
