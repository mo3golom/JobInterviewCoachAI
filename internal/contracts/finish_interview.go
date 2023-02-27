package contracts

import (
	"context"
	"github.com/google/uuid"
)

type FinishInterviewUseCase interface {
	FinishInterview(ctx context.Context, userID uuid.UUID) error
}
