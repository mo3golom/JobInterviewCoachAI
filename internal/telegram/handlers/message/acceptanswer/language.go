package acceptanswer

import (
	"job-interviewer/pkg/language"
)

const (
	textKeyProcessingAnswer language.TextKey = iota
	textKeyFinishInterview
	textKeyContinueInterview
)

func configLanguage() language.Storage {
	return language.NewLangStorage(
		map[language.Language]language.WordStorage{
			language.Russian: language.NewWordStorage(
				map[language.TextKey]string{
					textKeyFinishInterview:   "️⏏️️ Завершить",
					textKeyContinueInterview: "➡️ Продолжить",
				},
			),
			language.English: language.NewWordStorage(
				map[language.TextKey]string{
					textKeyProcessingAnswer: "Processing your answer...",
				},
			),
		},
	)
}
