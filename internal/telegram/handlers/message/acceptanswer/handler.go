package acceptanswer

import (
	"context"
	"errors"
	"fmt"
	interviewerContracts "job-interviewer/internal/interviewer/contracts"
	"job-interviewer/internal/telegram/handlers"
	"job-interviewer/internal/telegram/service"
	"job-interviewer/pkg/telegram"
	"job-interviewer/pkg/telegram/model"
	"job-interviewer/pkg/telegram/service/keyboard"
)

type Handler struct {
	acceptAnswerUC  interviewerContracts.AcceptAnswerUseCase
	keyboardService keyboard.Service
	service         service.Service
}

func NewHandler(
	auc interviewerContracts.AcceptAnswerUseCase,
	s service.Service,
	ks keyboard.Service,
) *Handler {
	return &Handler{
		acceptAnswerUC:  auc,
		service:         s,
		keyboardService: ks,
	}
}

func (h *Handler) Handle(ctx context.Context, request *model.Request, sender telegram.Sender) error {
	if request.Message == nil {
		return nil
	}

	answerMessageID, err := sender.Send(
		model.NewResponse(request.Chat.ID).SetText(handlers.ProcessingAnswerText),
	)
	response := model.NewResponse(request.Chat.ID)

	out, err := h.acceptAnswerUC.AcceptAnswer(
		ctx,
		interviewerContracts.AcceptAnswerIn{
			Answer: request.Message.Text,
			UserID: request.User.OriginalID,
		},
	)
	if errors.Is(err, interviewerContracts.ErrNextQuestionEmpty) {
		return h.service.FinishInterview(ctx, request, sender)
	}
	if errors.Is(err, interviewerContracts.ErrEmptyActiveInterview) {
		return nil //TODO: добавить возврат стартового меню
	}
	if err != nil {
		return err
	}

	inlineKeyboard, err := h.keyboardService.BuildInlineKeyboardInlineList(keyboard.BuildInlineKeyboardIn{
		Buttons: acceptAnswerButtons,
	})
	if err != nil {
		return err
	}
	err = sender.Update(
		answerMessageID,
		response.
			SetText(
				fmt.Sprintf(handlers.RobotPrefixText, out),
			).
			SetInlineKeyboardMarkup(inlineKeyboard),
	)
	if err != nil {
		return err
	}

	return nil
}
