package contracts

import (
	"context"
	"github.com/google/uuid"
	"job-interviewer/internal/interviewer/model"
)

type GetInterviewUseCase interface {
	GetActiveInterview(ctx context.Context, userID uuid.UUID) (*model.Interview, error)
}
