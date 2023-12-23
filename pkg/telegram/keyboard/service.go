package keyboard

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"job-interviewer/pkg/telegram/model"
)

const (
	maxCountButtonsWithoutChunks = 3
)

func BuildInlineKeyboardGrid(in BuildInlineKeyboardIn) (*tgbotapi.InlineKeyboardMarkup, error) {
	keyboard := make([][]tgbotapi.InlineKeyboardButton, 0, maxCountButtonsWithoutChunks)
	chunkSize := maxCountButtonsWithoutChunks
	for {
		if len(in.Buttons) == 0 {
			break
		}

		if len(in.Buttons) < chunkSize {
			chunkSize = len(in.Buttons)
		}

		row, err := buildInlineRow(in.Command, in.Buttons[:chunkSize])
		if err != nil {
			return nil, err
		}
		keyboard = append(keyboard, row)
		in.Buttons = in.Buttons[chunkSize:]
	}

	return &tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: keyboard,
	}, nil
}

func BuildInlineKeyboardList(in BuildInlineKeyboardIn) (*tgbotapi.InlineKeyboardMarkup, error) {
	keyboard := make([][]tgbotapi.InlineKeyboardButton, 0, len(in.Buttons))
	for _, btn := range in.Buttons {
		row, err := buildInlineRow(in.Command, []InlineButton{btn})
		if err != nil {
			return nil, err
		}

		keyboard = append(keyboard, row)

	}
	return &tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: keyboard,
	}, nil
}

func BuildInlineKeyboardInlineList(in BuildInlineKeyboardIn) (*tgbotapi.InlineKeyboardMarkup, error) {
	row, err := buildInlineRow(in.Command, in.Buttons)
	if err != nil {
		return nil, err
	}

	return &tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{row},
	}, nil
}

func BuildKeyboardGrid(in BuildKeyboardIn) *tgbotapi.ReplyKeyboardMarkup {
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

func BuildKeyboardCustomGrid(in BuildKeyboardCustomIn) *tgbotapi.ReplyKeyboardMarkup {
	keyboard := make([][]tgbotapi.KeyboardButton, 0, len(in.Buttons))
	for _, row := range in.Buttons {
		keyboard = append(
			keyboard,
			buildRow(row),
		)
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

func buildInlineRow(command *string, buttons []InlineButton) ([]tgbotapi.InlineKeyboardButton, error) {
	result := make([]tgbotapi.InlineKeyboardButton, 0, len(buttons))
	for _, button := range buttons {
		inlineButton, err := buildInlineButton(command, button)
		if err != nil {
			return nil, err
		}

		result = append(result, inlineButton)
	}

	return result, nil
}

func buildButton(button Button) tgbotapi.KeyboardButton {
	return tgbotapi.NewKeyboardButton(button.Value)
}

func buildInlineButton(command *string, button InlineButton) (tgbotapi.InlineKeyboardButton, error) {
	data := model.DataToString(button.Data)
	switch button.Type {
	case ButtonUrl:
		return tgbotapi.NewInlineKeyboardButtonURL(button.Value, data), nil
	}

	if command != nil {
		data = model.CommandWithDataToString(*command, button.Data)
	}

	if len(data) > 21 {
		return tgbotapi.InlineKeyboardButton{},
			fmt.Errorf("max data length in button = 21, '%s' = %d", data, len(data))
	}

	return tgbotapi.NewInlineKeyboardButtonData(button.Value, data), nil
}
