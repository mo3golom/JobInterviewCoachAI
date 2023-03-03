package keyboard

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

const (
	ButtonData ButtonType = iota
	ButtonUrl
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
		BuildKeyboardGrid(in BuildKeyboardIn) *tgbotapi.ReplyKeyboardMarkup
		BuildInlineKeyboardGrid(in BuildInlineKeyboardIn) *tgbotapi.InlineKeyboardMarkup
		BuildInlineKeyboardList(in BuildInlineKeyboardIn) *tgbotapi.InlineKeyboardMarkup
	}
)
