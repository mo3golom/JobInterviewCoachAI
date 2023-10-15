package acceptanswer

import (
	"job-interviewer/pkg/language"
)

const (
	textKeyProcessingAnswer          language.TextKey = "textKeyProcessingAnswer"
	textKeyVoiceMessageIsUnsupported language.TextKey = "textKeyVoiceMessageIsUnsupported"
)

func configLanguage() language.Storage {
	return language.NewLangStorage(
		map[language.Language]language.WordStorage{
			language.Russian: language.NewWordStorage(
				map[language.TextKey]string{
					textKeyVoiceMessageIsUnsupported: "Голосовые сообщения пока не поддерживаются, но мы работаем над этим",
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
