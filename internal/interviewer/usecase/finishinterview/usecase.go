package finishinterview

import (
	"context"
	"github.com/google/uuid"
	"job-interviewer/internal/interviewer/service/interview"
)

type UseCase struct {
	interviewService interview.Service
}

func NewUseCase(i interview.Service) *UseCase {
	return &UseCase{interviewService: i}
}

func (u *UseCase) FinishInterview(ctx context.Context, userID uuid.UUID) error {
	activeInterview, err := u.interviewService.FindActiveInterview(ctx, userID)
	if err != nil {
		return err
	}

	return u.interviewService.FinishInterview(ctx, activeInterview)
}
