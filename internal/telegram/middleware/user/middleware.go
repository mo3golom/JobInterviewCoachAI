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
	tgUser := request.User

	originalID, err := m.userUC.CreateOrGetUserToTelegram(ctx, tgUser.ID)
	if err != nil {
		return err
	}

	request.User.OriginalID = originalID
	return nil
}
