package user

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"job-interviewer/internal/interviewer/model"
	"job-interviewer/pkg/language"
	"job-interviewer/pkg/transactional"
)

type sqlxUser struct {
	ID   uuid.UUID         `db:"id"`
	Lang language.Language `db:"lang"`
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
		INTO "user" (id, lang) 
		VALUES (:id, :lang)
		ON CONFLICT DO NOTHING 
    `

	_, err := tx.NamedExecContext(
		ctx,
		query,
		sqlxUser{
			ID:   user.ID,
			Lang: user.Lang,
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

func (s *DefaultStorage) FindUserIDByTelegramID(ctx context.Context, tx transactional.Tx, telegramID int64) (*model.User, error) {
	query := `
		SELECT u.id, u.lang 
		FROM user_telegram as ut
		JOIN "user" as u on u.id = ut.user_id
		WHERE telegram_id = $1
    `

	var results []sqlxUser
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

	return &model.User{
		ID:   results[0].ID,
		Lang: results[0].Lang,
	}, nil
}

func (s *DefaultStorage) UpdateLanguage(ctx context.Context, tx transactional.Tx, userID uuid.UUID, language language.Language) error {
	query := `
       UPDATE "user" SET
        lang = :lang
       WHERE id = :id
    `

	in := sqlxUser{
		ID:   userID,
		Lang: language,
	}
	_, err := tx.NamedExecContext(ctx, query, in)
	return err
}
