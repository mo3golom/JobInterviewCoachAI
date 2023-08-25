package startinterview

import (
	"job-interviewer/pkg/language"
)

const (
	textKeyStartInterview language.TextKey = 1000
)

func configLanguage() language.Storage {
	return language.NewLangStorage(
		map[language.Language]language.WordStorage{
			language.Russian: language.NewWordStorage(
				map[language.TextKey]string{
					language.TextKey(QuestionContinueActiveInterview): "У вас есть уже активное интервью %s! Хотите продолжить?",
					language.TextKey(QuestionJobPosition):             "Выбери позицию, для которой хочешь пройти интервью:",
					textKeyStartInterview:                             "🆕 Начать новое интервью",
				},
			),
		},
	)
}
