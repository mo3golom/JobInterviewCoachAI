package about

import (
	"job-interviewer/pkg/language"
)

const (
	textKeyAbout        language.TextKey = "textKeyAbout"
	textKeyAboutCommand language.TextKey = "textKeyAboutCommand"
)

func configLanguage() language.Storage {
	return language.NewLangStorage(
		map[language.Language]language.WordStorage{
			language.Russian: language.NewWordStorage(
				map[language.TextKey]string{
					textKeyAbout: `Привет! Это интервью бот🚀

Я помогу вам попрактиковаться в собеседовании на английском языке. А после собеседования дам подробную обратную связь. 
Для этого я использую технологии от openai.

Я еще совсем маленький, поэтому иногда могу совершать ошибки 😰
Если у вас есть вопросы и предложения по поводу моей работы - напишите моему создателю %s
        `,
					textKeyAboutCommand: "🤖 О боте",
				},
			),
		},
	)
}
