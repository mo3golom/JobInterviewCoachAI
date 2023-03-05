package start

import (
	"job-interviewer/internal/telegram/handlers/command"
	"job-interviewer/pkg/telegram/service/keyboard"
)

var (
	buttons = []keyboard.Button{
		{
			Value: command.StartInterviewCommand,
		},
		{
			Value: "о боте",
		},
	}
)
