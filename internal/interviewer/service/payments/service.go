package payments

import (
	"context"
	"github.com/google/uuid"
	"job-interviewer/pkg/payments"
	externalSubscription "job-interviewer/pkg/subscription"
	"job-interviewer/pkg/transactional"
)

const (
	subscriptionPaymentType        payments.Type = "job_interviewer_subscription"
	subscriptionPaymentDescription               = "Оплата подписки на бот job interviewer coach ai"
	subscriptionPrice                            = 99
)

type DefaultService struct {
	paymentsService payments.Service
}

func NewDefaultService(
	paymentsService payments.Service,
	externalSubscriptionService externalSubscription.Service,
	transactionalTemplate transactional.Template,
) (*DefaultService, error) {
	err := paymentsService.RegisterPaymentPaidHandler(
		subscriptionPaymentType,
		&defaultPaidHandler{
			externalSubscriptionService: externalSubscriptionService,
			transactionalTemplate:       transactionalTemplate,
		},
	)
	if err != nil {
		return nil, err
	}

	return &DefaultService{
		paymentsService: paymentsService,
	}, nil
}

func (s *DefaultService) CreatePaymentForSubscription(ctx context.Context, userID uuid.UUID) (*CreatePaymentForSubscriptionOut, error) {
	payment, err := s.paymentsService.CreatePayment(
		ctx,
		&payments.CreatePaymentIn{
			IDK:         uuid.New(),
			UserID:      userID,
			Type:        subscriptionPaymentType,
			Description: subscriptionPaymentDescription,
			Amount:      subscriptionPrice,
		},
	)
	if err != nil {
		return nil, err
	}

	return &CreatePaymentForSubscriptionOut{
		RedirectURL: payment.RedirectURL,
	}, nil
}
