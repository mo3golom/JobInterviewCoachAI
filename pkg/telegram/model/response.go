package model

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type (
	Response struct {
		text                 string
		chatID               int64
		inlineKeyboardMarkup *tgbotapi.InlineKeyboardMarkup
		keyboardMarkup       *tgbotapi.ReplyKeyboardMarkup
	}
)

func NewResponse() Response {
	return Response{}
}

func (r Response) SetChatID(chatID int64) Response {
	r.chatID = chatID
	return r
}

func (r Response) SetText(value string) Response {
	r.text = value
	return r
}

func (r Response) SetInlineKeyboardMarkup(value *tgbotapi.InlineKeyboardMarkup) Response {
	r.inlineKeyboardMarkup = value
	return r
}

func (r Response) SetKeyboardMarkup(value *tgbotapi.ReplyKeyboardMarkup) Response {
	r.keyboardMarkup = value
	return r
}

func (r Response) ToMessageConfig() *tgbotapi.MessageConfig {
	result := tgbotapi.NewMessage(r.chatID, r.text)
	result.ParseMode = tgbotapi.ModeHTML
	result.ReplyMarkup = r.inlineKeyboardMarkup
	if r.keyboardMarkup != nil {
		result.ReplyMarkup = r.keyboardMarkup
	}

	return &result
}

func (r Response) ToEditMessageTextConfig() *tgbotapi.EditMessageTextConfig {
	result := tgbotapi.NewEditMessageText(r.chatID, 0, r.text)
	result.ParseMode = tgbotapi.ModeHTML
	result.ReplyMarkup = r.inlineKeyboardMarkup

	return &result
}
