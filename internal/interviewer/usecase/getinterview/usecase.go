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

func NewUseCase(
	interviewService interview.Service,
) *UseCase {
	return &UseCase{
		interviewService: interviewService,
	}
}

func (u *UseCase) FindActiveInterview(ctx context.Context, userID uuid.UUID) (*model.Interview, error) {
	return u.interviewService.FindActiveInterview(ctx, userID)
}

func (u *UseCase) GetAvailableValues() *model.InterviewAvailableValues {
	return &model.InterviewConfig
}
