package service

import (
	"job-interviewer/internal/telegram/handlers/command"
	"job-interviewer/pkg/telegram/service/keyboard"
)

var (
	getNextQuestionButtons = []keyboard.InlineButton{
		{
			Value: "#️⃣️ Завершить",
			Data:  []string{command.FinishInterviewCommand},
			Type:  keyboard.ButtonData,
		},
		{
			Value: "⏭️ Пропустить",
			Data:  []string{command.MarkQuestionAsSkip},
			Type:  keyboard.ButtonData,
		},
	}
)
