package payments

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"job-interviewer/pkg/helper"
	"job-interviewer/pkg/payments/gateway"
	"job-interviewer/pkg/payments/gateway/ym"
	"job-interviewer/pkg/payments/model"
	"job-interviewer/pkg/transactional"
	"net/http"
)

type (
	paymentCompletedHandlerItem struct {
		paymentType model.Type
		handler     PaymentCompletedHandler
	}

	DefaultService struct {
		repository            repository
		gateway               gateway.Gateway
		transactionalTemplate transactional.Template

		paidHandler     []paymentCompletedHandlerItem
		canceledHandler []paymentCompletedHandlerItem
	}
)

func NewPaymentsService(
	db *sqlx.DB,
	transactionalTemplate transactional.Template,
	YMShopID int64,
	YMSecretKey string,
) *DefaultService {
	return &DefaultService{
		transactionalTemplate: transactionalTemplate,
		repository: &DefaultRepository{
			db: db,
		},
		gateway: ym.NewYMGateway(
			YMShopID,
			YMSecretKey,
			&http.Client{},
		),
	}
}

func (s *DefaultService) CreatePayment(ctx context.Context, in *CreatePaymentIn) (*model.Payment, error) {
	existedPayment, err := s.repository.GetActivePaymentByUserID(ctx, in.UserID)
	if err != nil && !errors.Is(err, model.ErrPaymentNotFound) {
		return nil, err
	}

	if existedPayment == nil {
		existedPayment = &model.Payment{
			ID:          uuid.New(),
			UserID:      in.UserID,
			Amount:      in.Amount,
			Type:        in.Type,
			Description: in.Description,
		}
		err = s.transactionalTemplate.Execute(ctx, func(tx transactional.Tx) error {
			return s.repository.CreatePayment(ctx, tx, existedPayment)
		})
		if err != nil {
			return nil, err
		}
	}

	if existedPayment.ExternalID == nil {
		result, err := s.gateway.CreatePayment(
			ctx,
			&gateway.CreatePaymentIn{
				IDK:         in.IDK,
				Description: existedPayment.Description,
				Amount:      existedPayment.Amount.Normalize(),
			},
		)
		if err != nil {
			return nil, err
		}

		payment := &model.Payment{
			ID:          existedPayment.ID,
			Amount:      existedPayment.Amount,
			Type:        existedPayment.Type,
			Description: existedPayment.Description,
			Status:      model.StatusPending,

			ExternalID:  &result.ExternalID,
			RedirectURL: &result.RedirectURl,
		}
		err = s.transactionalTemplate.Execute(ctx, func(tx transactional.Tx) error {
			return s.repository.UpdatePayment(ctx, tx, payment)
		})
		if err != nil {
			return nil, err
		}

		return payment, nil
	}

	return existedPayment, nil
}

func (s *DefaultService) CheckPendingPayment(ctx context.Context, userID uuid.UUID) error {
	_, err := s.CheckPendingPaymentWithResult(ctx, userID)
	return err
}

func (s *DefaultService) CheckPendingPaymentWithResult(ctx context.Context, userID uuid.UUID) (*PaymentResult, error) {
	result := &PaymentResult{}

	payment, err := s.repository.GetActivePaymentByUserID(ctx, userID)
	if errors.Is(err, model.ErrPaymentNotFound) {
		return result, nil
	}
	if err != nil {
		return nil, err
	}

	result.PaymentType = payment.Type
	if payment.ExternalID == nil {
		return result, nil
	}
	paymentStatus, err := s.gateway.GetPaymentStatus(ctx, *payment.ExternalID)
	if err != nil {
		return nil, err
	}
	if payment.Status == *paymentStatus {
		return result, nil
	}

	payment.Status = *paymentStatus
	err = s.transactionalTemplate.Execute(ctx, func(tx transactional.Tx) error {
		return s.repository.UpdatePayment(ctx, tx, payment)
	})
	if err != nil {
		return nil, err
	}

	err = s.handlePaymentStatus(ctx, payment)
	if err != nil {
		return nil, err
	}

	result.Paid = *paymentStatus == model.StatusPaid
	return result, nil
}

func (s *DefaultService) handlePaymentStatus(ctx context.Context, payment *model.Payment) error {
	switch payment.Status {
	case model.StatusPaid:
		for _, item := range s.paidHandler {
			if item.paymentType != payment.Type {
				continue
			}

			err := item.handler.Handle(ctx, payment.UserID)
			if err == nil {
				continue
			}

			return err
		}
	case model.StatusCanceled:
		for _, item := range s.canceledHandler {
			if item.paymentType != payment.Type {
				continue
			}

			err := item.handler.Handle(ctx, payment.UserID)
			if err == nil {
				continue
			}

			return err
		}
	}

	return nil
}

func (s *DefaultService) RegisterPaymentPaidHandler(paymentType model.Type, handler PaymentCompletedHandler) error {
	handlers := helper.CopySlice(s.paidHandler)
	handlers = append(handlers, paymentCompletedHandlerItem{
		paymentType: paymentType,
		handler:     handler,
	})
	s.paidHandler = handlers
	return nil
}

func (s *DefaultService) RegisterPaymentCanceledHandler(paymentType model.Type, handler PaymentCompletedHandler) error {
	handlers := helper.CopySlice(s.canceledHandler)
	handlers = append(handlers, paymentCompletedHandlerItem{
		paymentType: paymentType,
		handler:     handler,
	})
	s.paidHandler = handlers
	return nil
}
