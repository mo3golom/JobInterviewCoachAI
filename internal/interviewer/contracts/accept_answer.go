package contracts

import (
	"context"
	"github.com/google/uuid"
	"job-interviewer/internal/interviewer/model"
)

type AcceptAnswerIn struct {
	UserID uuid.UUID
	Answer string
}

type AcceptAnswerUseCase interface {
	AcceptAnswer(ctx context.Context, in AcceptAnswerIn) error
	SkipQuestion(ctx context.Context, userID uuid.UUID) error
	GetAnswerSuggestion(ctx context.Context, userID uuid.UUID) (*model.AnswerSuggestion, error)
}
