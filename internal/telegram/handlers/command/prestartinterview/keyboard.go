package prestartinterview

import (
	"job-interviewer/internal/telegram/handlers/command"
	"job-interviewer/pkg/telegram/service/keyboard"
)

var (
	existsActiveInterviewButtons = []keyboard.InlineButton{
		{
			Value: "➡️ Продолжить",
			Data:  []string{command.GetNextQuestionCommand},
			Type:  keyboard.ButtonData,
		},
		{
			Value: "🆕 Новое",
			Data:  []string{command.ForceStartInterviewCommand},
			Type:  keyboard.ButtonData,
		},
	}
)
