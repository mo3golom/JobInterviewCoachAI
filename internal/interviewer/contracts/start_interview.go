package contracts

import (
	"context"
	"github.com/google/uuid"
)

type StartInterviewIn struct {
	UserID         uuid.UUID
	JobPosition    string
	QuestionsCount int64
}

type StartInterviewUseCase interface {
	StartInterview(ctx context.Context, in StartInterviewIn) error
}
