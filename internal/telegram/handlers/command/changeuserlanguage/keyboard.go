package changeuserlanguage

import (
	"job-interviewer/internal/telegram/language"
	"job-interviewer/pkg/telegram/service/keyboard"
)

var (
	buttons = []keyboard.InlineButton{
		{
			Value: "🇷🇺 Русский",
			Data:  []string{string(language.Russian)},
		},
		{
			Value: "🇺🇸 Английский",
			Data:  []string{string(language.English)},
		},
	}
)
