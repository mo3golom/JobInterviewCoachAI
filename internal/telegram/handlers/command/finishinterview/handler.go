package finishinterview

import (
	"context"
	"job-interviewer/internal/telegram/handlers/command"
	"job-interviewer/internal/telegram/service"
	"job-interviewer/pkg/telegram"
	"job-interviewer/pkg/telegram/model"
)

type Handler struct {
	service service.Service
}

func NewHandler(
	s service.Service,
) *Handler {
	return &Handler{
		service: s,
	}
}

func (h *Handler) Handle(ctx context.Context, request *model.Request, sender telegram.Sender) error {
	return h.service.FinishInterview(ctx, request, sender)
}

func (h *Handler) Command() string {
	return command.FinishInterviewCommand
}

func (h *Handler) Aliases() []string {
	return nil
}
