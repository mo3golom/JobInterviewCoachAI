package user

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"job-interviewer/internal/interviewer/model"
	"job-interviewer/pkg/language"
	"job-interviewer/pkg/transactional"
)

var (
	ErrEmptyUserResult = errors.New("empty user result")
)

type (
	TelegramUser struct {
		TelegramID int64
		Username   string
		FirstName  string
		LatName    string
	}

	Storage interface {
		CreateUser(ctx context.Context, tx transactional.Tx, user *model.User) error
		CreateTelegramToUser(ctx context.Context, tx transactional.Tx, in TelegramUser, userID uuid.UUID) error
		FindUserByTelegramID(ctx context.Context, tx transactional.Tx, telegramID int64) (*model.User, error)
		FindUserByID(ctx context.Context, userID uuid.UUID) (*model.User, error)
		UpdateLanguage(ctx context.Context, tx transactional.Tx, userID uuid.UUID, language language.Language) error
	}
)
