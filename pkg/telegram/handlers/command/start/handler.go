package start

import (
	"context"
	"job-interviewer/pkg/telegram"
	"job-interviewer/pkg/telegram/model"
	"job-interviewer/pkg/telegram/service/keyboard"
)

type Handler struct {
	keyboardService keyboard.Service
}

func NewHandler(
	k keyboard.Service,

) *Handler {
	return &Handler{
		keyboardService: k,
	}
}

func (h *Handler) Handle(_ context.Context, request *model.Request, sender telegram.Sender) error {
	_, err := sender.Send(
		model.NewResponse(request.Chat.ID).
			SetText("Hi").
			SetKeyboardMarkup(
				h.keyboardService.BuildKeyboardGrid(
					keyboard.BuildKeyboardIn{
						Buttons: buttons,
					},
				),
			),
	)

	return err
}

func (h *Handler) Command() string {
	return "/start"
}
