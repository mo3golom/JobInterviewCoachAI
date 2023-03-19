package user

import (
	"context"
	"job-interviewer/internal/interviewer/contracts"
	tgModel "job-interviewer/pkg/telegram/model"
)

type Middleware struct {
	userUC contracts.UserUseCase
}

func NewMiddleware(u contracts.UserUseCase) *Middleware {
	return &Middleware{
		userUC: u,
	}
}

func (m *Middleware) Handle(ctx context.Context, request *tgModel.Request) error {
	user, err := m.userUC.CreateOrGetUserToTelegram(
		ctx,
		&contracts.TgUserIn{
			ID:   request.User.ID,
			Lang: request.User.Lang,
		},
	)
	if err != nil {
		return err
	}

	request.User.OriginalID = user.ID
	if user.Lang != "" {
		request.User.Lang = user.Lang
	}
	return nil
}
