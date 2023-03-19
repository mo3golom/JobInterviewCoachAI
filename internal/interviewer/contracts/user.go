package contracts

import (
	"context"
	"github.com/google/uuid"
)

type (
	User struct {
		ID   uuid.UUID
		Lang string
	}

	TgUserIn struct {
		ID   int64
		Lang string
	}

	UserUseCase interface {
		CreateOrGetUserToTelegram(ctx context.Context, in *TgUserIn) (*User, error)
	}
)
