package prestartinterview

import "job-interviewer/pkg/telegram/service/keyboard"

var (
	existsActiveInterviewButtons = []keyboard.InlineButton{
		{
			Value: "➡️ Продолжить",
			Data:  []string{continueInterviewCommand},
			Type:  keyboard.ButtonData,
		},
		{
			Value: "🆕 Новое",
			Data:  []string{newInterviewCommand},
			Type:  keyboard.ButtonData,
		},
	}
)
