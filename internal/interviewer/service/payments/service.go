package payments

import (
	"context"
	"errors"
	"github.com/google/uuid"
	job_interviewer "job-interviewer"
	"job-interviewer/pkg/payments"
	"job-interviewer/pkg/payments/model"
	externalSubscription "job-interviewer/pkg/subscription"
	"job-interviewer/pkg/transactional"
	"job-interviewer/pkg/variables"
)

const (
	subscriptionPaymentType        model.Type = "job_interviewer_subscription"
	subscriptionPaymentDescription string     = "Оплата подписки на бот job interviewer coach ai"
)

type DefaultService struct {
	paymentsService payments.Service
	variables       variables.Repository
}

func NewDefaultService(
	paymentsService payments.Service,
	externalSubscriptionService externalSubscription.Service,
	transactionalTemplate transactional.Template,
	variables variables.Repository,
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
		variables:       variables,
	}, nil
}

func (s *DefaultService) CreatePaymentForSubscription(ctx context.Context, userID uuid.UUID) (*CreatePaymentForSubscriptionOut, error) {
	subscriptionPrice := s.variables.GetInt64(job_interviewer.MonthlySubscriptionPrice) * 100

	payment, err := s.paymentsService.CreatePayment(
		ctx,
		&payments.CreatePaymentIn{
			IDK:         uuid.New(),
			UserID:      userID,
			Type:        subscriptionPaymentType,
			Description: subscriptionPaymentDescription,
			Amount:      model.Penny(subscriptionPrice),
		},
	)
	if err != nil {
		return nil, err
	}
	if payment.RedirectURL == nil {
		return nil, errors.New("redirect url not found")
	}

	return &CreatePaymentForSubscriptionOut{
		RedirectURL: *payment.RedirectURL,
	}, nil
}
