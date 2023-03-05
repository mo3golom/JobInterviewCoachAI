package contracts

import (
	"context"
	"github.com/google/uuid"
	"job-interviewer/internal/interviewer/model"
)

type StartInterviewIn struct {
	UserID         uuid.UUID
	JobPosition    string
	JobLevel       model.JobLevel
	QuestionsCount int64
}

type StartInterviewUseCase interface {
	StartInterview(ctx context.Context, in StartInterviewIn) error
}
