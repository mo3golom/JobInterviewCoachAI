package contracts

import (
	"context"
	"github.com/google/uuid"
)

type UpdateQuestionUseCase interface {
	MarkActiveQuestionAsBad(ctx context.Context, userID uuid.UUID) error
	MarkActiveQuestionAsSkip(ctx context.Context, userID uuid.UUID) error
}
