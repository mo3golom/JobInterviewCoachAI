package start

import (
	"context"
	"job-interviewer/internal/telegram/service"
	"job-interviewer/pkg/telegram"
	"job-interviewer/pkg/telegram/model"
)

type Handler struct {
	service service.Service
}

func NewHandler(
	service service.Service,
) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Handle(_ context.Context, request *model.Request, sender telegram.Sender) error {
	return h.service.Start(request, sender)
}

func (h *Handler) Command() string {
	return "/start"
}
