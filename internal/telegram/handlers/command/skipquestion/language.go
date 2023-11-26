package skipquestion

import (
	"job-interviewer/pkg/language"
)

const (
	textKeyNoActiveInterview language.TextKey = "textKeyNoActiveInterview"
)

func configLanguage() language.Storage {
	return language.NewLangStorage(
		map[language.Language]language.WordStorage{
			language.Russian: language.NewWordStorage(
				map[language.TextKey]string{
					textKeyNoActiveInterview: "чтобы пропустить вопрос, у вас должна быть активная тренировка",
				},
			),
		},
	)
}
