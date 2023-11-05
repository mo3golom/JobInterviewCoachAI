package contracts

import (
	"context"
	"github.com/google/uuid"
	"job-interviewer/internal/interviewer/model"
)

type (
	GetInterviewUsecase interface {
		FindActiveInterview(ctx context.Context, userID uuid.UUID) (*model.Interview, error)
		GetAvailableValues() *model.InterviewAvailableValues
	}
)
