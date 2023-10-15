package checkpayment

import (
	"job-interviewer/pkg/language"
)

const (
	textKeyPaymentSuccess language.TextKey = "textKeyPaymentSuccess"
)

func configLanguage() language.Storage {
	return language.NewLangStorage(
		map[language.Language]language.WordStorage{
			language.Russian: language.NewWordStorage(
				map[language.TextKey]string{
					textKeyPaymentSuccess: "Подписка успешно оплачена, в ближайшее время она станет активной!",
				},
			),
		},
	)
}
