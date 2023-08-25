package service

import (
	"job-interviewer/pkg/language"
)

const (
	textKeyNotFoundActiveInterview language.TextKey = iota
	textKeyStartInterview
	textKeyFinishInterview
)

func configLanguage() language.Storage {
	return language.NewLangStorage(
		map[language.Language]language.WordStorage{
			language.Russian: language.NewWordStorage(
				map[language.TextKey]string{
					textKeyNotFoundActiveInterview: "I can`t find an active interview T-T",
					textKeyStartInterview:          "🆕 Начать новое интервью",
					textKeyFinishInterview:         "Interview’s over! Well done!",
				},
			),
		},
	)
}
