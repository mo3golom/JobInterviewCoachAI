package contracts

import (
	"context"
	"github.com/google/uuid"
	"job-interviewer/pkg/language"
)

type (
	User struct {
		ID   uuid.UUID
		Lang language.Language
	}

	TgUserIn struct {
		ID   int64
		Lang language.Language
	}

	UserUseCase interface {
		CreateOrGetUserToTelegram(ctx context.Context, in *TgUserIn) (*User, error)
		ChangeLanguage(ctx context.Context, userID uuid.UUID, language language.Language) error
	}
)
