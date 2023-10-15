package subscription

import (
	"context"
	"github.com/google/uuid"
	job_interviewer "job-interviewer"
	"job-interviewer/internal/interviewer/storage/messages"
	"job-interviewer/pkg/subscription"
	"job-interviewer/pkg/transactional"
	"job-interviewer/pkg/variables"
)

type DefaultService struct {
	subscriptionService subscription.Service
	messagesStorage     messages.Storage
	variables           variables.Repository
}

func NewService(
	subscriptionService subscription.Service,
	messagesStorage messages.Storage,
	variables variables.Repository,
) Service {
	return &DefaultService{
		subscriptionService: subscriptionService,
		messagesStorage:     messagesStorage,
		variables:           variables,
	}
}

func (s *DefaultService) IsAvailable(ctx context.Context, userID uuid.UUID) (*IsAvailableOut, error) {
	if !s.isEnabled() {
		return &IsAvailableOut{
			Result: true,
		}, nil
	}

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
	additionalChecker := &freeNextQuestionChecker{
		messagesStorage:    s.messagesStorage,
		userID:             userID,
		freeQuestionsCount: s.variables.GetInt64(job_interviewer.FreeQuestionsCount),
	}

	if !s.isEnabled() {
		result, err := additionalChecker.Check(ctx)
		if err != nil {
			return nil, err
		}

		return &IsAvailableOut{
			Result: result.Result,
			Reason: result.Reason,
		}, nil
	}

	result, err := s.subscriptionService.IsAvailable(
		ctx,
		userID,
		additionalChecker,
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
	if !s.isEnabled() {
		return nil
	}

	return s.subscriptionService.DecreaseFreeAttempts(ctx, tx, userID)
}

func (s *DefaultService) isEnabled() bool {
	return s.variables.GetBool(job_interviewer.PaidModelEnable)
}
