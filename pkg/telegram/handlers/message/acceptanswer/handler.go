package acceptanswer

import (
	"context"
	"errors"
	"fmt"
	"job-interviewer/internal/contracts"
	"job-interviewer/pkg/telegram"
	"job-interviewer/pkg/telegram/handlers"
	"job-interviewer/pkg/telegram/model"
	"job-interviewer/pkg/telegram/service/keyboard"
)

type Handler struct {
	getNextQuestionUC contracts.GetNextQuestionUseCase
	acceptAnswerUC    contracts.AcceptAnswerUseCase
	keyboardService   keyboard.Service
	finishHandler     telegram.Handler
}

func NewHandler(
	guc contracts.GetNextQuestionUseCase,
	auc contracts.AcceptAnswerUseCase,
	finishInterviewHandler telegram.Handler,
	ks keyboard.Service,
) *Handler {
	return &Handler{
		getNextQuestionUC: guc,
		acceptAnswerUC:    auc,
		finishHandler:     finishInterviewHandler,
		keyboardService:   ks,
	}
}

func (h *Handler) Handle(ctx context.Context, request *model.Request, sender telegram.Sender) error {
	if request.Message == nil {
		return nil
	}

	return h.acceptAnswer(ctx, request, sender)
}

func (h *Handler) acceptAnswer(ctx context.Context, request *model.Request, sender telegram.Sender) error {
	messageID, err := sender.Send(
		model.NewResponse(request.Chat.ID).SetText(handlers.ProcessingAnswerText),
	)
	response := model.NewResponse(request.Chat.ID)

	out, err := h.acceptAnswerUC.AcceptAnswer(
		ctx,
		contracts.AcceptAnswerIn{
			Answer: request.Message.Text,
			UserID: request.User.OriginalID,
		},
	)
	if errors.Is(err, contracts.ErrNextQuestionEmpty) {
		return h.finishHandler.Handle(ctx, request, sender)
	}
	if errors.Is(err, contracts.ErrEmptyActiveInterview) {
		return nil //TODO: добавить возврат стартового меню
	}
	if err != nil {
		return err
	}

	err = sender.Update(
		messageID,
		response.
			SetText(
				fmt.Sprintf(handlers.RobotPrefixText, out),
			).
			SetInlineKeyboardMarkup(
				h.keyboardService.BuildInlineKeyboardList(keyboard.BuildInlineKeyboardIn{
					Buttons: acceptAnswerButtons,
				}),
			),
	)
	if err != nil {
		return err
	}

	return nil
}
