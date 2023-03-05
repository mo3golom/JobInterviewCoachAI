package service

import (
	"context"
	"job-interviewer/pkg/telegram"
	"job-interviewer/pkg/telegram/model"
)

type Service interface {
	FinishInterview(ctx context.Context, request *model.Request, sender telegram.Sender) error
	GetNextQuestion(ctx context.Context, request *model.Request, sender telegram.Sender) error
	HideInlineKeyboard(request *model.Request, sender telegram.Sender) error
}
