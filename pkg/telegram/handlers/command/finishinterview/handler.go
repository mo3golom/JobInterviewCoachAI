package finishinterview

import (
	"context"
	"job-interviewer/internal/contracts"
	"job-interviewer/pkg/telegram"
	"job-interviewer/pkg/telegram/handlers"
	"job-interviewer/pkg/telegram/handlers/command"
	"job-interviewer/pkg/telegram/model"
)

type Handler struct {
	finishInterviewUC contracts.FinishInterviewUseCase
}

func NewHandler(
	f contracts.FinishInterviewUseCase,
) *Handler {
	return &Handler{
		finishInterviewUC: f,
	}
}

func (h *Handler) Handle(ctx context.Context, request *model.Request, sender telegram.Sender) error {
	err := h.finishInterviewUC.FinishInterview(ctx, request.User.OriginalID)
	if err != nil {
		return err
	}

	_, err = sender.Send(
		model.NewResponse(request.Chat.ID).
			SetText(handlers.FinishText),
	)
	return err
}

func (h *Handler) Command() string {
	return command.FinishInterviewCommand
}
