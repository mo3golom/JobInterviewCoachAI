package service

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"job-interviewer/pkg/language"
	"job-interviewer/pkg/telegram"
	"job-interviewer/pkg/telegram/model"
)

type Service interface {
	Start(request *model.Request, sender telegram.Sender) error
	FinishInterview(ctx context.Context, request *model.Request, sender telegram.Sender) error
	GetNextQuestion(ctx context.Context, request *model.Request, sender telegram.Sender) error
	SaveBotLastMessageID(ctx context.Context, chatID int64, lastBotMessageID int64) error
	HideInlineKeyboardForBotLastMessage(ctx context.Context, request *model.Request, sender telegram.Sender) error
	GetUserMainKeyboard(lang language.Language) *tgbotapi.ReplyKeyboardMarkup
}
