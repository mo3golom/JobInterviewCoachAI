package acceptanswer

import (
	"job-interviewer/internal/telegram/handlers/command"
	"job-interviewer/pkg/telegram/service/keyboard"
)

var (
	acceptAnswerButtons = []keyboard.InlineButton{
		{
			Value: "#️⃣️ Завершить",
			Data:  []string{command.FinishInterviewCommand},
			Type:  keyboard.ButtonData,
		},
		{
			Value: "▶️ Продолжить",
			Data:  []string{command.GetNextQuestionCommand},
			Type:  keyboard.ButtonData,
		},
	}
)
