package language

import "job-interviewer/pkg/language"

type (
	Dictionary interface {
		GetTexts() map[language.TextKey]string
	}

	Service interface {
		GetText(lang language.Language, key language.TextKey) string
		GetTextFromAllLanguages(key language.TextKey) []string
	}
)
