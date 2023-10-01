package startinterview

import (
	"context"
	"github.com/google/uuid"
	"job-interviewer/internal/interviewer/contracts"
	"job-interviewer/internal/interviewer/flow"
	"job-interviewer/internal/interviewer/service/interview"
	"job-interviewer/internal/interviewer/service/subscription"
)

type UseCase struct {
	interviewService    interview.Service
	interviewFlow       flow.InterviewFlow
	subscriptionService subscription.Service
}

func NewUseCase(
	i interview.Service,
	interviewFlow flow.InterviewFlow,
	subscriptionService subscription.Service,
) *UseCase {
	return &UseCase{
		interviewService:    i,
		interviewFlow:       interviewFlow,
		subscriptionService: subscriptionService,
	}
}

func (u *UseCase) StartInterview(ctx context.Context, in contracts.StartInterviewIn) error {
	available, err := u.subscriptionService.IsAvailable(ctx, in.UserID)
	if err != nil {
		return err
	}
	if !available.Result {
		return available.Reason
	}

	return u.interviewFlow.StartInterview(
		ctx,
		flow.StartInterviewIn{
			UserID:      in.UserID,
			JobPosition: in.Questions.JobPosition,
		},
	)
}

func (u *UseCase) ContinueInterview(ctx context.Context, userID uuid.UUID) error {
	available, err := u.subscriptionService.IsAvailable(ctx, userID)
	if err != nil {
		return err
	}
	if !available.Result {
		return available.Reason
	}

	return u.interviewFlow.ContinueInterview(ctx, userID)
}
