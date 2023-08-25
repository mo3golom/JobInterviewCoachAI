package finishinterview

import (
	"context"
	"github.com/google/uuid"
	"job-interviewer/internal/interviewer/flow"
)

type UseCase struct {
	interviewFlow flow.InterviewFlow
}

func NewUseCase(interviewFlow flow.InterviewFlow) *UseCase {
	return &UseCase{
		interviewFlow: interviewFlow,
	}
}

func (u *UseCase) FinishInterview(ctx context.Context, userID uuid.UUID) (string, error) {
	return u.interviewFlow.FinishInterview(ctx, userID)
}
