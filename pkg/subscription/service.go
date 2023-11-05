package subscription

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"job-interviewer/internal/interviewer/contracts"
	"job-interviewer/pkg/transactional"
	"time"
)

type DefaultService struct {
	storage           storage
	freeAttemptsCount int64
}

func NewSubscriptionService(
	db *sqlx.DB,
	freeAttemptsCount int64,
) Service {
	return &DefaultService{
		storage:           &defaultStorage{db: db},
		freeAttemptsCount: freeAttemptsCount,
	}
}

func (d *DefaultService) CreateUser(ctx context.Context, tx transactional.Tx, userID uuid.UUID) error {
	return d.storage.createUser(
		ctx,
		tx,
		userID,
		UserTypeFree,
		d.freeAttemptsCount,
	)
}

func (d *DefaultService) IsAvailable(ctx context.Context, userID uuid.UUID, additionalCheckers ...Checker) (*IsAvailableOut, error) {
	userInfo, err := d.storage.getUserInfo(ctx, userID)
	if err != nil {
		return nil, err
	}

	switch userInfo.Type {
	case UserTypeFree:
		if userInfo.FreeAttempts != nil && *userInfo.FreeAttempts <= 0 {
			return &IsAvailableOut{
				Result: false,
				Reason: contracts.ErrFreeAttemptsHaveExpired,
			}, nil
		}

		return isAvailableAdditional(ctx, UserTypeFree, additionalCheckers)
	case UserTypePaid:
		if userInfo.Subscription == nil || !userInfo.Subscription.IsActive() {
			return &IsAvailableOut{
				Result: false,
				Reason: contracts.ErrPaidSubscriptionHasExpired,
			}, nil
		}

		return isAvailableAdditional(ctx, UserTypePaid, additionalCheckers)
	}

	return &IsAvailableOut{
		Result: false,
	}, nil

}

func (d *DefaultService) DecreaseFreeAttempts(ctx context.Context, tx transactional.Tx, userID uuid.UUID) error {
	userInfo, err := d.storage.getUserInfo(ctx, userID)
	if err != nil {
		return err
	}

	if userInfo.Type == UserTypePaid {
		return nil
	}

	if userInfo.FreeAttempts != nil && *userInfo.FreeAttempts <= 0 {
		return nil
	}

	return d.storage.decreaseFreeAttempts(ctx, tx, userID)

}

func (d *DefaultService) ActivateSubscription(ctx context.Context, tx transactional.Tx, userID uuid.UUID, plan Plan) error {
	switch plan {
	case PlanMonth:
		return d.storage.upsertSubscription(ctx, tx, userID, time.Now(), time.Now().AddDate(0, 1, 0))
	}

	return nil
}

func isAvailableAdditional(ctx context.Context, userType UserType, checkers []Checker) (*IsAvailableOut, error) {
	if len(checkers) == 0 {
		return &IsAvailableOut{
			Result: true,
		}, nil
	}

	for _, checker := range checkers {
		if checker.Type() != userType {
			continue
		}

		isAvailableResult, err := checker.Check(ctx)
		if err != nil {
			return nil, err
		}
		if isAvailableResult.Result {
			continue
		}

		return isAvailableResult, nil
	}

	return &IsAvailableOut{
		Result: true,
	}, nil
}
