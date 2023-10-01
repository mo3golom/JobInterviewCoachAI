package service

import (
	"job-interviewer/pkg/language"
)

const (
	textKeyNotFoundActiveInterview language.TextKey = iota
	textKeyStartInterview
	textKeyFinishInterview
	textKeyFreeQuestionsIsEnd
	textKeySubscriptionHasExpired
	textKeySubscribe
)

func configLanguage() language.Storage {
	return language.NewLangStorage(
		map[language.Language]language.WordStorage{
			language.Russian: language.NewWordStorage(
				map[language.TextKey]string{
					textKeyNotFoundActiveInterview: "I can`t find an active interview T-T",
					textKeyStartInterview:          "🚀 Начать новое интервью",
					textKeyFinishInterview:         "Interview’s over! Well done!",
					textKeyFreeQuestionsIsEnd:      "В бесплатном плане можно ответить лишь на огранниченное число вопросов. Завершаю интервью!",
					textKeySubscriptionHasExpired:  "Чтобы продолжить работу с ботом, необходимо оплатить подписку",
					textKeySubscribe: `Спасибо что проявили интерес к нашему боту! К сожалению ваша подписка истекла или число бесплатных попыток закончилось :( 
Чтобы продолжить пользоваться ботом приобретите подписку всего за 99р в месяц :)`,
				},
			),
		},
	)
}
