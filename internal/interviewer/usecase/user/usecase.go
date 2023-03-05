package user

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"job-interviewer/internal/interviewer/model"
	"job-interviewer/internal/interviewer/storage/user"
	"job-interviewer/pkg/transactional"
)

type UseCase struct {
	userStorage           user.Storage
	transactionalTemplate transactional.Template
}

func NewUseCase(u user.Storage, t transactional.Template) *UseCase {
	return &UseCase{
		userStorage:           u,
		transactionalTemplate: t,
	}
}

func (u *UseCase) CreateOrGetUserToTelegram(ctx context.Context, tgUserID int64) (uuid.UUID, error) {
	var originalID uuid.UUID
	err := u.transactionalTemplate.Execute(ctx, func(tx transactional.Tx) error {
		existed, err := u.userStorage.FindUserIDByTelegramID(ctx, tx, tgUserID)
		if existed != nil {
			originalID = *existed
			return nil
		}

		if err != nil && !errors.Is(err, user.ErrEmptyUserResult) {
			return err
		}

		originalID = uuid.New()
		newUser := &model.User{
			ID: originalID,
		}
		err = u.userStorage.CreateUser(ctx, tx, newUser)
		if err != nil {
			return err
		}

		return u.userStorage.CreateTelegramToUser(ctx, tx, tgUserID, newUser.ID)
	})

	return originalID, err
}
