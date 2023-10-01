package subscription

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"job-interviewer/pkg/transactional"
	"time"
)

type (
	sqlxPaidSubscription struct {
		UserID  uuid.UUID `db:"user_id"`
		StartAt time.Time `db:"start_at"`
		EndAt   time.Time `db:"end_at"`
	}

	defaultStorage struct {
		db *sqlx.DB
	}
)

func (d *defaultStorage) createUser(ctx context.Context, tx transactional.Tx, userID uuid.UUID, userType UserType, freeAttempts int64) error {
	const userQuery = `
		INSERT INTO subscription_users (
			 user_id,
			 type
		) VALUES (
		   $1,
		   $2	  
		) ON CONFLICT (user_id) DO NOTHING 
	`

	_, err := tx.ExecContext(
		ctx,
		userQuery,
		userID,
		string(userType),
	)
	if err != nil {
		return err
	}

	const freeAttemptsQuery = `
		INSERT INTO subscription_free_attempt (
			 user_id,
			 attempts
		) VALUES (
		   $1,
		   $2	  
		) ON CONFLICT (user_id) DO NOTHING 
	`

	_, err = tx.ExecContext(
		ctx,
		freeAttemptsQuery,
		userID,
		freeAttempts,
	)
	return err
}

func (d *defaultStorage) decreaseFreeAttempts(ctx context.Context, tx transactional.Tx, userID uuid.UUID) error {
	const query = `
       UPDATE subscription_free_attempt SET
        attempts = attempts - 1,
        updated_at = now()   
       WHERE user_id = $1
    `

	_, err := tx.ExecContext(ctx, query, userID)
	return err
}

func (d *defaultStorage) getUserInfo(ctx context.Context, userID uuid.UUID) (*Info, error) {
	const query = `
		SELECT su.type, s.start_at, s.end_at ,sfa.attempts  FROM subscription_users as su
		LEFT JOIN subscription as s on s.user_id = su.user_id
		LEFT JOIN subscription_free_attempt as sfa on sfa.user_id = su.user_id
		WHERE su.user_id = $1
	
`

	var results []struct {
		Type     string     `db:"type"`
		StartAt  *time.Time `db:"start_at"`
		EndAt    *time.Time `db:"end_at"`
		Attempts *int64     `db:"attempts"`
	}

	err := d.db.SelectContext(
		ctx,
		&results,
		query,
		userID,
	)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, ErrEmptyUserInfo
	}

	var subscription *Subscription
	if results[0].StartAt != nil && results[0].EndAt != nil {
		subscription = &Subscription{
			StartAt: *results[0].StartAt,
			EndAt:   *results[0].EndAt,
		}
	}

	return &Info{
		Type:         UserType(results[0].Type),
		FreeAttempts: results[0].Attempts,
		Subscription: subscription,
	}, nil
}

func (d *defaultStorage) upsertSubscription(ctx context.Context, tx transactional.Tx, userID uuid.UUID, startAt time.Time, endAt time.Time) error {
	const userQuery = `
		INSERT INTO subscription_users (
			 user_id,
			 type
		) VALUES (
		   $1,
		   $2	  
		) ON CONFLICT (user_id) DO UPDATE SET type = $2
	`
	_, err := tx.ExecContext(
		ctx,
		userQuery,
		userID,
		string(UserTypePaid),
	)
	if err != nil {
		return err
	}

	const subscriptionQuery = `
        INSERT 
		INTO subscription (user_id, start_at, end_at) 
		VALUES (:user_id, :start_at, :end_at)
		ON CONFLICT (user_id) DO UPDATE SET
		    start_at = :start_at, 
		    end_at = :end_at,
		    updated_at = now()
    `
	_, err = tx.NamedExecContext(
		ctx,
		subscriptionQuery,
		sqlxPaidSubscription{
			UserID:  userID,
			StartAt: startAt,
			EndAt:   endAt,
		},
	)
	return err
}
