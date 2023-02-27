package keyboard

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

const (
	ButtonUrl ButtonType = iota
	ButtonData
)

type (
	ButtonType int64

	Button struct {
		Value string
		Type  ButtonType
	}

	InlineButton struct {
		Value string
		Data  []string
		Type  ButtonType
	}

	BuildInlineKeyboardIn struct {
		Command *string
		Buttons []InlineButton
	}

	BuildKeyboardIn struct {
		Buttons []Button
	}

	Service interface {
		BuildKeyboard(in BuildKeyboardIn) *tgbotapi.ReplyKeyboardMarkup
		BuildInlineKeyboard(in BuildInlineKeyboardIn) *tgbotapi.InlineKeyboardMarkup
	}
)
