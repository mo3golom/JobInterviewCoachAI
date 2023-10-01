package errors

import (
	"context"
	"errors"
	"job-interviewer/internal/interviewer/contracts"
	"job-interviewer/internal/telegram/service"
	"job-interviewer/pkg/logger"
	"job-interviewer/pkg/telegram"
	"job-interviewer/pkg/telegram/model"
)

type SubscribeErrorHandler struct {
	service service.Service
	logger  logger.Logger
}

func NewSubscribeErrorHandler(
	service service.Service,
	logger logger.Logger,
) telegram.ErrorHandler {
	return &SubscribeErrorHandler{
		service: service,
		logger:  logger,
	}
}

func (e *SubscribeErrorHandler) Handle(_ context.Context, err error, _ *model.Request, sender telegram.Sender) bool {
	if errors.Is(err, contracts.ErrFreeAttemptsHaveExpired) || errors.Is(err, contracts.ErrPaidSubscriptionHasExpired) {
		localErr := e.service.ShowSubscribeMessage(sender)
		if localErr != nil {
			e.logger.Error("subscribe error handler error", localErr)
		}

		return true
	}

	return false
}
