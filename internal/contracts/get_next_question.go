package contracts

import (
	"context"
	"github.com/google/uuid"
	"job-interviewer/internal/model"
)

type GetNextQuestionUseCase interface {
	GetNextQuestion(ctx context.Context, userID uuid.UUID) (*model.Question, error)
}
