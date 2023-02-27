package tguser

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"job-interviewer/internal/model"
	"job-interviewer/internal/storage/user"
	tgModel "job-interviewer/pkg/telegram/model"
	"job-interviewer/pkg/transactional"
)

type Middleware struct {
	userStorage           user.Storage
	transactionalTemplate transactional.Template
}

func NewMiddleware(u user.Storage, t transactional.Template) *Middleware {
	return &Middleware{
		userStorage:           u,
		transactionalTemplate: t,
	}
}

func (m *Middleware) Handle(ctx context.Context, request *tgModel.Request) error {
	tgUser := request.User

	return m.transactionalTemplate.Execute(ctx, func(tx transactional.Tx) error {
		existed, err := m.userStorage.FindUserIDByTelegramID(ctx, tx, int64(tgUser.ID))
		if existed != nil {
			request.User.OriginalID = *existed
		}
		if err != nil && !errors.Is(err, user.ErrEmptyUserResult) {
			return err
		}

		newUser := &model.User{
			ID: uuid.New(),
		}
		err = m.userStorage.CreateUser(ctx, tx, newUser)
		if err != nil {
			return err
		}

		return m.userStorage.CreateTelegramToUser(ctx, tx, int64(tgUser.ID), newUser.ID)
	})

}
