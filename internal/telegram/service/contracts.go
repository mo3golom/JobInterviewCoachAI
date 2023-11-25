package service

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"job-interviewer/pkg/language"
	"job-interviewer/pkg/telegram"
	"job-interviewer/pkg/telegram/model"
)

type Service interface {
	FinishInterview(ctx context.Context, request *model.Request, sender telegram.Sender, updateMessageID ...int64) error
	GetNextQuestion(ctx context.Context, request *model.Request, sender telegram.Sender, updateMessageID ...int64) error
	GetUserMainKeyboard(lang language.Language) *tgbotapi.ReplyKeyboardMarkup
	ShowSubscribeMessage(sender telegram.Sender) error
}
