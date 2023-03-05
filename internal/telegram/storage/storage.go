package storage

import (
	"context"
	"github.com/jmoiron/sqlx"
)

type (
	sqlxTelegramBotDetails struct {
		ChatID           int64 `db:"chat_id"`
		LastBotMessageID int64 `db:"last_bot_message_id"`
	}

	DefaultStorage struct {
		db *sqlx.DB
	}
)

func NewStorage(db *sqlx.DB) *DefaultStorage {
	return &DefaultStorage{db: db}
}

func (s *DefaultStorage) GetBotLastMessageID(ctx context.Context, chatID int64) (int64, error) {
	query := `
		SELECT tbd.last_bot_message_id
		FROM telegram_bot_details as tbd
		WHERE tbd.chat_id = $1 
		LIMIT 1
    `

	var results []sqlxTelegramBotDetails
	err := s.db.SelectContext(
		ctx,
		&results,
		query,
		chatID,
	)
	if err != nil {
		return 0, err
	}
	if len(results) == 0 {
		return 0, ErrEmptyTelegramBotDetailsResult
	}

	return results[0].LastBotMessageID, nil
}

func (s *DefaultStorage) UpsertTelegramBotDetails(ctx context.Context, in UpsertTelegramBotDetailsIn) error {
	query := `
        INSERT 
		INTO telegram_bot_details (chat_id, last_bot_message_id) 
		VALUES (:chat_id, :last_bot_message_id)
		ON CONFLICT (chat_id) DO UPDATE
        SET  last_bot_message_id = :last_bot_message_id, updated_at = now()
    `

	_, err := s.db.NamedExecContext(
		ctx,
		query,
		sqlxTelegramBotDetails{
			ChatID:           in.ChatID,
			LastBotMessageID: in.LastBotMessageID,
		},
	)
	return err
}
