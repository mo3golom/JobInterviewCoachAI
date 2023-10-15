package start

import (
	"job-interviewer/pkg/language"
)

const (
	textKeyStart language.TextKey = "textKeyStart"
)

func configLanguage() language.Storage {
	return language.NewLangStorage(
		map[language.Language]language.WordStorage{
			language.Russian: language.NewWordStorage(
				map[language.TextKey]string{
					textKeyStart: `Привет! Я ваш идеальный помощник в подготовке к собеседованию на разных языках!

Просто нажмите "🚀 Начать новое интервью", выберите специализацию, на которой хотите потренироваться. А дальше я сгенерирую список соответствующих вопросов, используя модель GPT от openAI. 

Как только вы ответите на вопрос, я дам вам конструктивную обратную связь, и буду готов отправить следующий вопрос :)
        `,
				},
			),
		},
	)
}
