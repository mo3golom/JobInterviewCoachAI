package acceptanswer

import (
	"job-interviewer/pkg/language"
)

const (
	textKeyProcessingAnswer language.TextKey = iota
)

func configLanguage() language.Storage {
	return language.NewLangStorage(
		map[language.Language]language.WordStorage{
			language.English: language.NewWordStorage(
				map[language.TextKey]string{
					textKeyProcessingAnswer: "Processing your answer...",
				},
			),
		},
	)
}
