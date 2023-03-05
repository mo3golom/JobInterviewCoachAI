package contracts

import (
	"context"
	"github.com/google/uuid"
	"job-interviewer/internal/interviewer/model"
)

type GetNextQuestionUseCase interface {
	GetNextQuestion(ctx context.Context, userID uuid.UUID) (*model.Question, error)
}
