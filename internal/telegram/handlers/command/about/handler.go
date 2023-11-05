package about

import (
	"context"
	"fmt"
	job_interviewer "job-interviewer"
	"job-interviewer/internal/telegram/handlers/command"
	"job-interviewer/pkg/language"
	"job-interviewer/pkg/telegram"
	"job-interviewer/pkg/telegram/model"
	"job-interviewer/pkg/variables"
)

type Handler struct {
	variablesRepository variables.Repository
	languageStorage     language.Storage
}

func NewHandler(
	variablesRepository variables.Repository,
) *Handler {
	return &Handler{
		variablesRepository: variablesRepository,
		languageStorage:     configLanguage(),
	}
}

func (h *Handler) Handle(_ context.Context, _ *model.Request, sender telegram.Sender) error {
	userLang := language.Russian

	contact := h.variablesRepository.GetString(job_interviewer.TGContact)
	_, err := sender.Send(
		model.NewResponse().
			SetText(
				fmt.Sprintf(h.languageStorage.GetText(userLang, textKeyAbout), contact),
			),
	)

	return err
}

func (h *Handler) Command() string {
	return command.AboutCommand
}

func (h *Handler) Aliases() []string {
	return []string{
		h.languageStorage.GetText(language.Russian, textKeyAboutCommand),
	}
}
