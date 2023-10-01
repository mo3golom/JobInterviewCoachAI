package subscription

import (
	"context"
	"github.com/google/uuid"
	"job-interviewer/pkg/transactional"
	"time"
)

type (
	IsAvailableOut struct {
		Result bool
		Reason error
	}

	Service interface {
		CreateUser(ctx context.Context, tx transactional.Tx, userID uuid.UUID) error
		IsAvailable(ctx context.Context, userID uuid.UUID, additionalCheckers ...Checker) (*IsAvailableOut, error)
		DecreaseFreeAttempts(ctx context.Context, tx transactional.Tx, userID uuid.UUID) error
		ActivateSubscription(ctx context.Context, tx transactional.Tx, userID uuid.UUID, plan Plan) error
	}

	Checker interface {
		Check(ctx context.Context) (*IsAvailableOut, error)
		Type() UserType
	}

	storage interface {
		createUser(ctx context.Context, tx transactional.Tx, userID uuid.UUID, userType UserType, freeAttempts int64) error
		decreaseFreeAttempts(ctx context.Context, tx transactional.Tx, userID uuid.UUID) error
		upsertSubscription(ctx context.Context, tx transactional.Tx, userID uuid.UUID, startAt time.Time, endAt time.Time) error
		getUserInfo(ctx context.Context, userID uuid.UUID) (*Info, error)
	}
)
