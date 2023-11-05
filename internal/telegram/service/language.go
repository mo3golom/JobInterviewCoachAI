package service

import (
	"job-interviewer/pkg/language"
)

const (
	textKeyNotFoundActiveInterview language.TextKey = "textKeyNotFoundActiveInterview"
	textKeyStartInterview          language.TextKey = "textKeyStartInterview"
	textKeyFinishInterview         language.TextKey = "textKeyFinishInterview"
	textKeyFreeQuestionsIsEnd      language.TextKey = "textKeyFreeQuestionsIsEnd"
	textKeySubscriptionHasExpired  language.TextKey = "textKeySubscriptionHasExpired"
	textKeySubscribe               language.TextKey = "textKeySubscribe"
	textKeyBuySubscription         language.TextKey = "textKeyBuySubscription"
	textKeyAbout                   language.TextKey = "textKeyAbout"
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
Месячная подписка стоит %d рублей в месяц. 
После активации подписки вам будет доступно неограниченное тренировок с неограниченным числом вопросов! Если вы готовы, нажмите кнопку "%s" :)`,
					textKeyBuySubscription: "💸 Купить подписку",
					textKeyAbout:           "🤖 О боте",
				},
			),
		},
	)
}
