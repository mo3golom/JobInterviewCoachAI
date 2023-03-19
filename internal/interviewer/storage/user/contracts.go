package user

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"job-interviewer/internal/interviewer/model"
	"job-interviewer/pkg/transactional"
)

var (
	ErrEmptyUserResult = errors.New("empty user result")
)

type Storage interface {
	CreateUser(ctx context.Context, tx transactional.Tx, user *model.User) error
	CreateTelegramToUser(ctx context.Context, tx transactional.Tx, telegramID int64, userID uuid.UUID) error
	FindUserIDByTelegramID(ctx context.Context, tx transactional.Tx, telegramID int64) (*model.User, error)
}
