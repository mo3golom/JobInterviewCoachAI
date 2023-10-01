package user

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"job-interviewer/internal/interviewer/contracts"
	"job-interviewer/internal/interviewer/model"
	"job-interviewer/internal/interviewer/storage/user"
	"job-interviewer/pkg/language"
	"job-interviewer/pkg/subscription"
	"job-interviewer/pkg/transactional"
)

type UseCase struct {
	userStorage           user.Storage
	transactionalTemplate transactional.Template
	subscription          subscription.Service
}

func NewUseCase(
	u user.Storage,
	t transactional.Template,
	subscription subscription.Service,
) *UseCase {
	return &UseCase{
		userStorage:           u,
		transactionalTemplate: t,
		subscription:          subscription,
	}
}

func (u *UseCase) CreateOrGetUserToTelegram(ctx context.Context, in *contracts.TgUserIn) (*contracts.User, error) {
	var originalUser *model.User

	err := u.transactionalTemplate.Execute(ctx, func(tx transactional.Tx) error {
		existed, err := u.userStorage.FindUserByTelegramID(ctx, tx, in.ID)
		if existed != nil {
			originalUser = existed

			return u.subscription.CreateUser(ctx, tx, originalUser.ID)
		}

		if err != nil && !errors.Is(err, user.ErrEmptyUserResult) {
			return err
		}

		originalUser = &model.User{
			ID:   uuid.New(),
			Lang: in.Lang,
		}
		err = u.userStorage.CreateUser(ctx, tx, originalUser)
		if err != nil {
			return err
		}

		err = u.subscription.CreateUser(ctx, tx, originalUser.ID)
		if err != nil {
			return err
		}

		return u.userStorage.CreateTelegramToUser(ctx, tx, in.ID, originalUser.ID)
	})
	if err != nil {
		return nil, err
	}

	return &contracts.User{
		ID:   originalUser.ID,
		Lang: originalUser.Lang,
	}, nil
}

func (u *UseCase) ChangeLanguage(ctx context.Context, userID uuid.UUID, language language.Language) error {
	return u.transactionalTemplate.Execute(ctx, func(tx transactional.Tx) error {
		return u.userStorage.UpdateLanguage(ctx, tx, userID, language)
	})
}
