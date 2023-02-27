package contracts

import (
	"context"
	"github.com/google/uuid"
	"job-interviewer/internal/model"
)

type StartInterviewIn struct {
	UserID         uuid.UUID
	JobPosition    string
	JobLevel       model.JobLevel
	QuestionsCount int64
}

type StartInterviewUsecase interface {
	StartInterview(ctx context.Context, in StartInterviewIn) error
}
