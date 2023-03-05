package storage

import (
	"context"
	"errors"
)

var (
	ErrEmptyTelegramBotDetailsResult = errors.New("not found telegram bot details")
)

type (
	UpsertTelegramBotDetailsIn struct {
		ChatID           int64
		LastBotMessageID int64
	}

	Storage interface {
		GetBotLastMessageID(ctx context.Context, chatID int64) (int64, error)
		UpsertTelegramBotDetails(ctx context.Context, in UpsertTelegramBotDetailsIn) error
	}
)
