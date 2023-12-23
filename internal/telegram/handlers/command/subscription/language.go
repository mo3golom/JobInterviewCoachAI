package subscription

import "job-interviewer/pkg/language"

const (
	textKeySubscription      = "subscription"
	textKeySubscriptionAbout = "subscription_about"
)

func configLanguage() language.Storage {
	return language.NewLangStorage(
		map[language.Language]language.WordStorage{
			language.Russian: language.NewWordStorage(
				map[language.TextKey]string{
					textKeySubscription: "⭐ Подписка",
					textKeySubscriptionAbout: `Сейчас бот находится в бета-тесте и абсолютно бесплатен, но каждое интервью ограничено 5 вопросами.

В будущем для бота станет доступна подписка всего за <b>99 рублей в месяц</b>, которая позволит проходить неограниченное число тренировок с неограниченным числом вопросов! 🚀

Также всем новым пользователям будет доступен бесплатный пробный режим, чтобы Вы сами могли оценить полезность бота для себя и принять решение о приобретении подписки 😊`,
				},
			),
		},
	)
}
