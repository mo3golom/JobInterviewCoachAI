package payments

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"job-interviewer/pkg/helper"
)

type DefaultService struct {
	repository repository
	gateway    Gateway

	paidHandler     map[Type]PaymentCompletedHandler
	canceledHandler map[Type]PaymentCompletedHandler
}

func (s *DefaultService) CreatePayment(ctx context.Context, in *CreatePaymentIn) (*Payment, error) {
	existedPayment, err := s.repository.GetPendingPaymentByUserID(ctx, in.UserID)
	if err != nil && !errors.Is(err, ErrPaymentNotFound) {
		return nil, err
	}
	if existedPayment != nil {
		return existedPayment, nil
	}

	newPayment := &NewPayment{
		ID:          uuid.New(),
		UserID:      in.UserID,
		Amount:      in.Amount,
		Type:        in.Type,
		Description: in.Description,
	}
	err = s.repository.CreatePayment(
		ctx,
		newPayment,
	)
	if err != nil {
		return nil, err
	}

	result, err := s.gateway.CreatePayment(
		ctx,
		&GatewayCreatePaymentIn{
			IDK:         in.IDK,
			Description: newPayment.Description,
			Amount:      newPayment.Amount.Normalize(),
		},
	)
	if err != nil {
		return nil, err
	}

	payment := &Payment{
		ID:          newPayment.ID,
		ExternalID:  result.ExternalID,
		Amount:      newPayment.Amount,
		Type:        newPayment.Type,
		Description: newPayment.Description,
		RedirectURL: result.RedirectURl,
		Status:      StatusPending,
	}
	err = s.repository.UpdatePayment(ctx, payment)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

func (s *DefaultService) CheckPendingPayment(ctx context.Context, userID uuid.UUID) error {
	payment, err := s.repository.GetPendingPaymentByUserID(ctx, userID)
	if errors.Is(err, ErrPaymentNotFound) {
		return nil
	}
	if err != nil {
		return err
	}

	paymentStatus, err := s.gateway.GetPaymentStatus(ctx, payment.ExternalID)
	if err != nil {
		return err
	}
	if payment.Status == *paymentStatus {
		return nil
	}

	payment.Status = *paymentStatus
	err = s.repository.UpdatePayment(ctx, payment)
	if err != nil {
		return err
	}

	return s.handlePaymentStatus(ctx, payment)
}

func (s *DefaultService) handlePaymentStatus(ctx context.Context, payment *Payment) error {
	switch payment.Status {
	case StatusPaid:
		handler, ok := s.paidHandler[payment.Type]
		if !ok {
			return nil
		}

		return handler.Handle(ctx, payment.UserID)
	case StatusCanceled:
		handler, ok := s.canceledHandler[payment.Type]
		if !ok {
			return nil
		}

		return handler.Handle(ctx, payment.UserID)
	}

	return nil
}

func (s *DefaultService) RegisterPaymentPaidHandler(paymentType Type, handler PaymentCompletedHandler) error {
	_, ok := s.paidHandler[paymentType]
	if ok {
		return ErrHandlerAlreadyRegistered
	}

	handlers := helper.CopyMap[Type, PaymentCompletedHandler](s.paidHandler)
	handlers[paymentType] = handler
	s.paidHandler = handlers
	return nil
}

func (s *DefaultService) RegisterPaymentCanceledHandler(paymentType Type, handler PaymentCompletedHandler) error {
	_, ok := s.canceledHandler[paymentType]
	if ok {
		return ErrHandlerAlreadyRegistered
	}

	handlers := helper.CopyMap[Type, PaymentCompletedHandler](s.canceledHandler)
	handlers[paymentType] = handler
	s.canceledHandler = handlers
	return nil
}
