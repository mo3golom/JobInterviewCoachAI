package getnextquestion

import (
	"context"
	"errors"
	"fmt"
	"job-interviewer/internal/contracts"
	"job-interviewer/pkg/telegram"
	"job-interviewer/pkg/telegram/handlers"
	"job-interviewer/pkg/telegram/handlers/command"
	"job-interviewer/pkg/telegram/model"
)

type Handler struct {
	getNextQuestionUC contracts.GetNextQuestionUseCase
	finishHandler     telegram.Handler
}

func NewHandler(
	g contracts.GetNextQuestionUseCase,
	finishInterviewHandler telegram.Handler,
) *Handler {
	return &Handler{
		getNextQuestionUC: g,
		finishHandler:     finishInterviewHandler,
	}
}

func (h *Handler) Handle(ctx context.Context, request *model.Request, sender telegram.Sender) error {
	response := model.NewResponse(request.Chat.ID)

	question, err := h.getNextQuestionUC.GetNextQuestion(ctx, request.User.OriginalID)
	if errors.Is(err, contracts.ErrNextQuestionEmpty) {
		return h.finishHandler.Handle(ctx, request, sender)
	}
	if errors.Is(err, contracts.ErrEmptyActiveInterview) {
		_, err = sender.Send(response.SetText(handlers.NoActiveInterviewText))
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}

	_, err = sender.Send(
		response.SetText(fmt.Sprintf(handlers.RobotPrefixText, question.Text)),
	)
	return err
}

func (h *Handler) Command() string {
	return command.GetNextQuestionCommand
}
