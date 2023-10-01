package service

import (
	"job-interviewer/internal/telegram/handlers/command"
	"job-interviewer/pkg/telegram/service/keyboard"
)

var (
	getNextQuestionButtons = []keyboard.InlineButton{
		{
			Value: "🚜 Я все!",
			Data:  []string{command.FinishInterviewCommand},
			Type:  keyboard.ButtonData,
		},
		{
			Value: "🙋 Подскажи",
			Data:  []string{command.GetAnswerSuggestionCommand},
			Type:  keyboard.ButtonData,
		},
		{
			Value: "💭 Не хочу отвечать",
			Data:  []string{command.SkipQuestionCommand},
			Type:  keyboard.ButtonData,
		},
	}

	subscribeButtons = []keyboard.InlineButton{
		{
			Value: "💸 Купить подписку",
			Data:  []string{command.FinishInterviewCommand},
			Type:  keyboard.ButtonData,
		},
	}
)
