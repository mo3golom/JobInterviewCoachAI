package subscription

import (
	"context"
	"job-interviewer/internal/telegram/handlers/command"
	"job-interviewer/pkg/language"
	"job-interviewer/pkg/telegram"
	"job-interviewer/pkg/telegram/model"
)

type Handler struct {
	languageStorage language.Storage
}

func NewHandler() *Handler {
	return &Handler{
		languageStorage: configLanguage(),
	}
}

func (h *Handler) Handle(_ context.Context, _ *model.Request, sender telegram.Sender) error {
	userLang := language.Russian

	_, err := sender.Send(
		model.NewResponse().
			SetText(h.languageStorage.GetText(userLang, textKeySubscriptionAbout)),
	)

	return err
}

func (h *Handler) Command() string {
	return command.SubscriptionCommand
}

func (h *Handler) Aliases() []string {
	return []string{
		h.languageStorage.GetText(language.Russian, textKeySubscription),
	}
}
