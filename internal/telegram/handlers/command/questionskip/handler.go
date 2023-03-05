package questionskip

import (
	"context"
	"job-interviewer/internal/interviewer/contracts"
	"job-interviewer/internal/telegram/handlers/command"
	"job-interviewer/internal/telegram/service"
	"job-interviewer/pkg/telegram"
	"job-interviewer/pkg/telegram/model"
)

type Handler struct {
	updateQuestionUC contracts.UpdateQuestionUseCase
	service          service.Service
}

func NewHandler(
	updateQuestionUC contracts.UpdateQuestionUseCase,
	service service.Service,
) *Handler {
	return &Handler{
		updateQuestionUC: updateQuestionUC,
		service:          service,
	}
}

func (h *Handler) Handle(ctx context.Context, request *model.Request, sender telegram.Sender) error {
	err := h.updateQuestionUC.MarkActiveQuestionAsSkip(ctx, request.User.OriginalID)
	if err != nil {
		return err
	}

	return h.service.GetNextQuestion(ctx, request, sender)
}

func (h *Handler) Command() string {
	return command.MarkQuestionAsSkip
}
