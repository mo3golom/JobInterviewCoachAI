package contracts

import (
	"context"
	"github.com/google/uuid"
)

type AcceptAnswerIn struct {
	UserID uuid.UUID
	Answer string
}

type AcceptAnswerUseCase interface {
	AcceptAnswer(ctx context.Context, in AcceptAnswerIn) (string, error)
}
