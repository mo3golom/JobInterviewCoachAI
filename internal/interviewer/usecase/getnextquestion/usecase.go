package getnextquestion

import (
	"context"
	"github.com/google/uuid"
	"job-interviewer/internal/interviewer/flow"
	"job-interviewer/internal/interviewer/model"
	"job-interviewer/internal/interviewer/service/subscription"
)

type UseCase struct {
	interviewFlow       flow.InterviewFlow
	subscriptionService subscription.Service
}

func NewUseCase(
	interviewFlow flow.InterviewFlow,
	subscriptionService subscription.Service,
) *UseCase {
	return &UseCase{
		interviewFlow:       interviewFlow,
		subscriptionService: subscriptionService,
	}
}

func (u *UseCase) GetNextQuestion(ctx context.Context, userID uuid.UUID) (*model.Question, error) {
	available, err := u.subscriptionService.IsAvailableNextQuestion(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !available.Result {
		return nil, available.Reason
	}

	return u.interviewFlow.NextQuestion(ctx, userID)
}
