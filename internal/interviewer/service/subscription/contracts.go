package subscription

import (
	"context"
	"github.com/google/uuid"
	"job-interviewer/pkg/transactional"
)

type (
	IsAvailableOut struct {
		Result bool
		Reason error
	}

	Service interface {
		IsAvailable(ctx context.Context, userID uuid.UUID) (*IsAvailableOut, error)
		IsAvailableNextQuestion(ctx context.Context, userID uuid.UUID) (*IsAvailableOut, error)
		DecreaseFreeAttempts(ctx context.Context, tx transactional.Tx, userID uuid.UUID) error
	}
)
