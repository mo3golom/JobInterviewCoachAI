package language

import "job-interviewer/pkg/language"

type (
	Dictionary interface {
		GetTexts() map[language.TextKey]string
	}

	Service interface {
		InitUserLanguage(lang language.Language) error
		GetUserLanguageText(key language.TextKey) string
		InitInterviewLanguage(lang language.Language) error
		GetInterviewLanguageText(key language.TextKey) string
	}
)
