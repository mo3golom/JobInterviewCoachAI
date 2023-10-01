package subscription

import (
	"context"
	"github.com/google/uuid"
	"job-interviewer/internal/interviewer/storage/messages"
	"job-interviewer/pkg/subscription"
	"job-interviewer/pkg/transactional"
)

const (
	defaultFreeQuestions = 1
)

type DefaultService struct {
	subscriptionService subscription.Service
	messagesStorage     messages.Storage
}

func NewService(
	subscriptionService subscription.Service,
	messagesStorage messages.Storage,
) Service {
	return &DefaultService{
		subscriptionService: subscriptionService,
		messagesStorage:     messagesStorage,
	}
}

func (s *DefaultService) IsAvailable(ctx context.Context, userID uuid.UUID) (*IsAvailableOut, error) {
	result, err := s.subscriptionService.IsAvailable(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &IsAvailableOut{
		Result: result.Result,
		Reason: result.Reason,
	}, nil
}

func (s *DefaultService) IsAvailableNextQuestion(ctx context.Context, userID uuid.UUID) (*IsAvailableOut, error) {
	result, err := s.subscriptionService.IsAvailable(
		ctx,
		userID,
		&freeNextQuestionChecker{
			messagesStorage: s.messagesStorage,
			userID:          userID,
		},
	)
	if err != nil {
		return nil, err
	}

	return &IsAvailableOut{
		Result: result.Result,
		Reason: result.Reason,
	}, nil
}

func (s *DefaultService) DecreaseFreeAttempts(ctx context.Context, tx transactional.Tx, userID uuid.UUID) error {
	return s.subscriptionService.DecreaseFreeAttempts(ctx, tx, userID)
}
