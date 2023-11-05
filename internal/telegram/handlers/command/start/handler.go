package start

import (
	"context"
	"job-interviewer/internal/telegram/service"
	"job-interviewer/pkg/language"
	"job-interviewer/pkg/telegram"
	"job-interviewer/pkg/telegram/model"
)

type Handler struct {
	service         service.Service
	languageStorage language.Storage
}

func NewHandler(
	service service.Service,
) *Handler {
	return &Handler{
		service:         service,
		languageStorage: configLanguage(),
	}
}

func (h *Handler) Handle(_ context.Context, _ *model.Request, sender telegram.Sender) error {
	userLang := language.Russian

	_, err := sender.Send(
		model.NewResponse().
			SetText(h.languageStorage.GetText(userLang, textKeyStart)).
			SetKeyboardMarkup(h.service.GetUserMainKeyboard(userLang)),
	)

	return err
}

func (h *Handler) Command() string {
	return "/start"
}

func (h *Handler) Aliases() []string {
	return nil
}
