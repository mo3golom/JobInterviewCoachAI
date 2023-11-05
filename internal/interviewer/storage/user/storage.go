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
	Username   string    `db:"username"`
	Firstname  string    `db:"first_name"`
	LastName   string    `db:"last_name"`
}

type DefaultStorage struct {
	db *sqlx.DB
}

func NewStorage(db *sqlx.DB) *DefaultStorage {
	return &DefaultStorage{db: db}
}

func (s *DefaultStorage) CreateUser(ctx context.Context, tx transactional.Tx, user *model.User) error {
	const query = `
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

func (s *DefaultStorage) CreateTelegramToUser(ctx context.Context, tx transactional.Tx, in TelegramUser, userID uuid.UUID) error {
	const query = `
		INSERT 
		INTO user_telegram (user_id, telegram_id, username, first_name, last_name) 
		VALUES (:user_id, :telegram_id, :username, :first_name, :last_name)
		ON CONFLICT DO NOTHING 
    `

	_, err := tx.NamedExecContext(
		ctx,
		query,
		sqlxTelegramUser{
			UserID:     userID,
			TelegramID: in.TelegramID,
			Username:   in.Username,
			Firstname:  in.FirstName,
			LastName:   in.LatName,
		},
	)
	return err
}

func (s *DefaultStorage) FindUserByTelegramID(ctx context.Context, tx transactional.Tx, telegramID int64) (*model.User, error) {
	const query = `
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

func (s *DefaultStorage) FindUserByID(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	const query = `
		SELECT u.id, u.lang
		FROM "user" as u
		WHERE u.id = $1
    `

	var results []sqlxUser
	err := s.db.SelectContext(
		ctx,
		&results,
		query,
		userID,
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
	const query = `
       UPDATE "user" SET
        lang = :lang,
        updated_at = now()
       WHERE id = :id
    `

	in := sqlxUser{
		ID:   userID,
		Lang: language,
	}
	_, err := tx.NamedExecContext(ctx, query, in)
	return err
}
