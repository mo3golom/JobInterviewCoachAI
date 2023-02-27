package keyboard

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"job-interviewer/pkg/telegram/model"
)

const (
	maxCountButtonsWithoutChunks = 3
)

type DefaultService struct {
}

func (s *DefaultService) BuildInlineKeyboard(in BuildInlineKeyboardIn) *tgbotapi.InlineKeyboardMarkup {
	keyboard := make([][]tgbotapi.InlineKeyboardButton, 0, maxCountButtonsWithoutChunks)
	chunkSize := maxCountButtonsWithoutChunks
	for {
		if len(in.Buttons) == 0 {
			break
		}

		if len(in.Buttons) < chunkSize {
			chunkSize = len(in.Buttons)
		}

		keyboard = append(
			keyboard,
			buildInlineRow(in.Command, in.Buttons[:chunkSize]),
		)
		in.Buttons = in.Buttons[chunkSize:]
	}

	return &tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: keyboard,
	}
}

func (s *DefaultService) BuildKeyboard(in BuildKeyboardIn) *tgbotapi.ReplyKeyboardMarkup {
	keyboard := make([][]tgbotapi.KeyboardButton, 0, maxCountButtonsWithoutChunks)
	chunkSize := maxCountButtonsWithoutChunks
	for {
		if len(in.Buttons) == 0 {
			break
		}

		if len(in.Buttons) < chunkSize {
			chunkSize = len(in.Buttons)
		}

		keyboard = append(
			keyboard,
			buildRow(in.Buttons[:chunkSize]),
		)
		in.Buttons = in.Buttons[chunkSize:]
	}

	return &tgbotapi.ReplyKeyboardMarkup{
		ResizeKeyboard: true,
		Keyboard:       keyboard,
	}
}

func buildRow(buttons []Button) []tgbotapi.KeyboardButton {
	result := make([]tgbotapi.KeyboardButton, 0, len(buttons))
	for _, button := range buttons {
		result = append(
			result,
			buildButton(button),
		)
	}

	return result
}

func buildInlineRow(command *string, buttons []InlineButton) []tgbotapi.InlineKeyboardButton {
	result := make([]tgbotapi.InlineKeyboardButton, 0, len(buttons))
	for _, button := range buttons {
		result = append(
			result,
			buildInlineButton(command, button),
		)
	}

	return result
}

func buildButton(button Button) tgbotapi.KeyboardButton {
	return tgbotapi.NewKeyboardButton(button.Value)
}

func buildInlineButton(command *string, button InlineButton) tgbotapi.InlineKeyboardButton {
	data := model.DataToString(button.Data)
	switch button.Type {
	case ButtonUrl:
		return tgbotapi.NewInlineKeyboardButtonURL(button.Value, data)
	}

	if command != nil {
		data = model.CommandWithDataToString(*command, button.Data)
	}
	return tgbotapi.NewInlineKeyboardButtonData(button.Value, data)
}
