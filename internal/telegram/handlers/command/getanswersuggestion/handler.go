package getanswersuggestion

import (
	"context"
	"errors"
	"fmt"
	interviewerContracts "job-interviewer/internal/interviewer/contracts"
	"job-interviewer/internal/telegram/handlers"
	"job-interviewer/internal/telegram/handlers/command"
	"job-interviewer/pkg/telegram"
	"job-interviewer/pkg/telegram/model"
)

type Handler struct {
	acceptAnswerUC interviewerContracts.AcceptAnswerUseCase
}

func NewHandler(
	auc interviewerContracts.AcceptAnswerUseCase,
) *Handler {
	return &Handler{
		acceptAnswerUC: auc,
	}
}

func (h *Handler) Handle(ctx context.Context, request *model.Request, sender telegram.Sender) error {
	messageID, err := sender.Send(
		model.
			NewResponse().
			SetText(fmt.Sprintf("%s %s", handlers.RobotPrefix, "loading suggestions...")),
	)

	result, err := h.acceptAnswerUC.GetAnswerSuggestion(ctx, request.User.OriginalID)
	if errors.Is(err, interviewerContracts.ErrEmptyActiveInterview) {
		return nil
	}
	if err != nil {
		return err
	}

	return sender.Update(
		messageID,
		model.
			NewResponse().
			SetText(fmt.Sprintf("%s %s", handlers.RobotPrefix, result.Text)),
	)
}

func (h *Handler) Command() string {
	return command.GetAnswerSuggestionCommand
}

func (h *Handler) Aliases() []string {
	return nil
}
