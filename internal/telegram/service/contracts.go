package service

import (
	"context"
	"job-interviewer/pkg/telegram"
	"job-interviewer/pkg/telegram/model"
)

type Service interface {
	Start(request *model.Request, sender telegram.Sender) error
	FinishInterview(ctx context.Context, request *model.Request, sender telegram.Sender) error
	GetNextQuestion(ctx context.Context, request *model.Request, sender telegram.Sender) error
	SaveBotLastMessageID(ctx context.Context, chatID int64, lastBotMessageID int64) error
	HideInlineKeyboardForBotLastMessage(ctx context.Context, request *model.Request, sender telegram.Sender) error
}
