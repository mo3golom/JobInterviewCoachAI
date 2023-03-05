package getinterview

import (
	"context"
	"github.com/google/uuid"
	"job-interviewer/internal/interviewer/model"
	"job-interviewer/internal/interviewer/service/interview"
)

type UseCase struct {
	interviewService interview.Service
}

func NewUseCase(i interview.Service) *UseCase {
	return &UseCase{interviewService: i}
}

func (u *UseCase) GetActiveInterview(ctx context.Context, userID uuid.UUID) (*model.Interview, error) {
	activeInterview, err := u.interviewService.FindActiveInterview(ctx, userID)
	if err != nil {
		return nil, err
	}

	return activeInterview, err
}
