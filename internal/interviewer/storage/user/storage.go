package user

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"job-interviewer/internal/interviewer/model"
	"job-interviewer/pkg/transactional"
)

type sqlxUser struct {
	ID uuid.UUID `db:"id"`
}

type sqlxTelegramUser struct {
	UserID     uuid.UUID `db:"user_id"`
	TelegramID int64     `db:"telegram_id"`
}

type DefaultStorage struct {
	db *sqlx.DB
}

func NewStorage(db *sqlx.DB) *DefaultStorage {
	return &DefaultStorage{db: db}
}

func (s *DefaultStorage) CreateUser(ctx context.Context, tx transactional.Tx, user *model.User) error {
	query := `
		INSERT 
		INTO "user" (id) 
		VALUES (:id)
		ON CONFLICT DO NOTHING 
    `

	_, err := tx.NamedExecContext(
		ctx,
		query,
		sqlxUser{
			ID: user.ID,
		},
	)
	return err
}

func (s *DefaultStorage) CreateTelegramToUser(ctx context.Context, tx transactional.Tx, telegramID int64, userID uuid.UUID) error {
	query := `
		INSERT 
		INTO user_telegram (user_id, telegram_id) 
		VALUES (:user_id, :telegram_id)
		ON CONFLICT DO NOTHING 
    `

	_, err := tx.NamedExecContext(
		ctx,
		query,
		sqlxTelegramUser{
			UserID:     userID,
			TelegramID: telegramID,
		},
	)
	return err
}

func (s *DefaultStorage) FindUserIDByTelegramID(ctx context.Context, tx transactional.Tx, telegramID int64) (*uuid.UUID, error) {
	query := `
		SELECT user_id 
		FROM user_telegram
		WHERE telegram_id = $1
    `

	var results []sqlxTelegramUser
	err := tx.SelectContext(
		ctx,
		&results,
		query,
		telegramID,
	)

	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, ErrEmptyUserResult
	}

	userID := results[0].UserID
	return &userID, nil
}
