package contracts

import (
	"context"
	"github.com/google/uuid"
)

type UserUseCase interface {
	CreateOrGetUserToTelegram(ctx context.Context, tgUserID int64) (uuid.UUID, error)
}
