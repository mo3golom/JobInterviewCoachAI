package service

import (
	"job-interviewer/pkg/language"
)

const (
	textKeyNotFoundActiveInterview = "textKeyNotFoundActiveInterview"
	textKeyStartInterview          = "textKeyStartInterview"
	textKeyFinishInterview         = "textKeyFinishInterview"
	textKeyFreeQuestionsIsEnd      = "textKeyFreeQuestionsIsEnd"
	textKeySubscriptionHasExpired  = "textKeySubscriptionHasExpired"
	textKeySubscribe               = "textKeySubscribe"
	textKeyBuySubscription         = "textKeyBuySubscription"
	textKeyAbout                   = "textKeyAbout"
	textKeyFinishNoActiveInterview = "textKeyFinishNoActiveInterview"
	textKeySubscription            = "subscription"
	textKeyImDone                  = "imDone"
	textKeySuggestion              = "suggestion"
	textKeySkipQuestion            = "skipQuestion"
)

func configLanguage() language.Storage {
	return language.NewLangStorage(
		map[language.Language]language.WordStorage{
			language.Russian: language.NewWordStorage(
				map[language.TextKey]string{
					textKeyNotFoundActiveInterview: "I can`t find an active interview 😥",
					textKeyStartInterview:          "🚀 Начать новое интервью",
					textKeyFinishInterview:         "Interview’s over! Well done!",
					textKeyFreeQuestionsIsEnd:      "В бесплатном плане можно ответить лишь на ограниченное число вопросов. Завершаю интервью!",
					textKeySubscriptionHasExpired:  "Чтобы продолжить работу с ботом, необходимо оплатить подписку",
					textKeySubscribe: `Спасибо что проявили интерес к нашему боту! К сожалению ваша подписка истекла или число бесплатных попыток закончилось 😔 
Месячная подписка стоит <b>%d</b> рублей в месяц. 
После активации подписки вам будет доступно неограниченное тренировок с неограниченным числом вопросов! Если вы готовы, нажмите кнопку <b>"%s"</b> 😉`,
					textKeySubscription:            "⭐ Подписка",
					textKeyBuySubscription:         "💸 Купить подписку",
					textKeyAbout:                   "🤖 О боте",
					textKeyFinishNoActiveInterview: "чтобы закончить тренировку, у вас должна быть активная тренировка",
					textKeyImDone:                  "🚜 Я все!",
					textKeySuggestion:              "🙋 Подскажи",
					textKeySkipQuestion:            "💭 Не хочу отвечать",
				},
			),
		},
	)
}
