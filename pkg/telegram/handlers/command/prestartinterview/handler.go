package prestartinterview

import (
	"context"
	"errors"
	"job-interviewer/internal/contracts"
	"job-interviewer/pkg/telegram"
	"job-interviewer/pkg/telegram/handlers/command"
	"job-interviewer/pkg/telegram/model"
	"job-interviewer/pkg/telegram/service/keyboard"
)

type Handler struct {
	getInterviewUC         contracts.GetInterviewUseCase
	keyboardService        keyboard.Service
	startHandler           telegram.Handler
	getNextQuestionHandler telegram.Handler
}

func NewHandler(
	k keyboard.Service,
	g contracts.GetInterviewUseCase,
	startInterviewHandler telegram.Handler,
	getNextQuestionHandler telegram.Handler,
) *Handler {
	return &Handler{
		keyboardService:        k,
		getInterviewUC:         g,
		startHandler:           startInterviewHandler,
		getNextQuestionHandler: getNextQuestionHandler,
	}
}

func (h *Handler) Handle(ctx context.Context, request *model.Request, sender telegram.Sender) error {
	if len(request.Data) == 1 {
		switch request.Data[0] {
		case newInterviewCommand:
			request.Data = []string{} // обнуляем данные
			return h.startHandler.Handle(ctx, request, sender)
		case continueInterviewCommand:
			return h.getNextQuestionHandler.Handle(ctx, request, sender)
		}

		return nil
	}

	_, err := h.getInterviewUC.GetActiveInterview(ctx, request.User.OriginalID)
	if errors.Is(err, contracts.ErrEmptyActiveInterview) {
		return h.startHandler.Handle(ctx, request, sender)
	}
	if err != nil {
		return err
	}

	currentCommand := h.Command()
	_, err = sender.Send(
		model.NewResponse(request.Chat.ID).
			SetText(activeInterviewExistsText).
			SetInlineKeyboardMarkup(
				h.keyboardService.BuildInlineKeyboardGrid(
					keyboard.BuildInlineKeyboardIn{
						Command: &currentCommand,
						Buttons: existsActiveInterviewButtons,
					},
				),
			),
	)
	return err
}

func (h *Handler) Command() string {
	return command.StartInterviewCommand
}
