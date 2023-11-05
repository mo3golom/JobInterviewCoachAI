package finishinterview

import (
	"context"
	"github.com/google/uuid"
	"job-interviewer/internal/interviewer/flow"
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

func (u *UseCase) FinishInterview(ctx context.Context, userID uuid.UUID) (string, error) {
	available, err := u.subscriptionService.IsAvailable(ctx, userID)
	if err != nil {
		return "", err
	}
	if !available.Result {
		return "", available.Reason
	}

	return u.interviewFlow.FinishInterview(ctx, userID)
}
