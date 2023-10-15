package paysubscription

import (
	"job-interviewer/pkg/language"
)

const (
	textKeyPayStart       language.TextKey = "textKeyPayStart"
	textKeyPayRedirectURL language.TextKey = "textKeyPayRedirectURL"
	textKeyPay            language.TextKey = "textKeyPay"
	textKeyCheckPayment   language.TextKey = "textKeyCheckPayment"
)

func configLanguage() language.Storage {
	return language.NewLangStorage(
		map[language.Language]language.WordStorage{
			language.Russian: language.NewWordStorage(
				map[language.TextKey]string{
					textKeyPayStart:       "💸 Купить подписку",
					textKeyPayRedirectURL: "Создана ссылка на оплату подписки. После оплаты нажмите на кнопку \"Проверить оплату\"",
					textKeyPay:            "Оплатить",
					textKeyCheckPayment:   "Проверить оплату",
				},
			),
		},
	)
}
