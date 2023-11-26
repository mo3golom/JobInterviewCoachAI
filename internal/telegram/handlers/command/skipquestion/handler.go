package skipquestion

import (
	"context"
	"errors"
	"fmt"
	interviewerContracts "job-interviewer/internal/interviewer/contracts"
	"job-interviewer/internal/telegram/handlers"
	"job-interviewer/internal/telegram/handlers/command"
	"job-interviewer/internal/telegram/service"
	"job-interviewer/pkg/language"
	"job-interviewer/pkg/telegram"
	"job-interviewer/pkg/telegram/model"
)

type Handler struct {
	acceptAnswerUC  interviewerContracts.AcceptAnswerUseCase
	service         service.Service
	languageStorage language.Storage
}

func NewHandler(
	auc interviewerContracts.AcceptAnswerUseCase,
	s service.Service,
) *Handler {
	return &Handler{
		acceptAnswerUC:  auc,
		service:         s,
		languageStorage: configLanguage(),
	}
}

func (h *Handler) Handle(ctx context.Context, request *model.Request, sender telegram.Sender) error {
	if request.CallbackID != nil {
		err := sender.SendCallback(*request.CallbackID)
		if err != nil {
			return err
		}
	}

	messageID, err := sender.Send(
		model.
			NewResponse().
			SetText(fmt.Sprintf("%s %s", handlers.RobotPrefix, "loading...")),
	)
	if err != nil {
		return err
	}

	err = h.acceptAnswerUC.SkipQuestion(ctx, request.User.OriginalID)
	if errors.Is(err, interviewerContracts.ErrEmptyActiveInterview) {
		return sender.Update(
			messageID,
			model.
				NewResponse().
				SetText(h.languageStorage.GetText(language.Russian, textKeyNoActiveInterview)),
		)
	}
	if err != nil {
		return err
	}

	return h.service.GetNextQuestion(ctx, request, sender, messageID)
}

func (h *Handler) Command() string {
	return command.SkipQuestionCommand
}

func (h *Handler) Aliases() []string {
	return nil
}
