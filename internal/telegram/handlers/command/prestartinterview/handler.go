package prestartinterview

import (
	"context"
	"errors"
	interviewContracts "job-interviewer/internal/interviewer/contracts"
	"job-interviewer/internal/telegram/handlers/command"
	languageService "job-interviewer/internal/telegram/language"
	"job-interviewer/internal/telegram/service"
	"job-interviewer/pkg/telegram"
	"job-interviewer/pkg/telegram/model"
	"job-interviewer/pkg/telegram/service/keyboard"
)

type Handler struct {
	getInterviewUC  interviewContracts.GetInterviewUseCase
	keyboardService keyboard.Service
	startHandler    telegram.Handler
	service         service.Service
	languageService languageService.Service
}

func NewHandler(
	k keyboard.Service,
	g interviewContracts.GetInterviewUseCase,
	startInterviewHandler telegram.Handler,
	s service.Service,
	l languageService.Service,
) *Handler {
	return &Handler{
		keyboardService: k,
		getInterviewUC:  g,
		startHandler:    startInterviewHandler,
		service:         s,
		languageService: l,
	}
}

func (h *Handler) Handle(ctx context.Context, request *model.Request, sender telegram.Sender) error {
	userLang := request.User.Lang

	_, err := h.getInterviewUC.GetActiveInterview(ctx, request.User.OriginalID)
	if errors.Is(err, interviewContracts.ErrEmptyActiveInterview) {
		return h.startHandler.Handle(ctx, request, sender)
	}
	if err != nil {
		return err
	}

	inlineKeyboard, err := h.keyboardService.BuildInlineKeyboardGrid(
		keyboard.BuildInlineKeyboardIn{
			Buttons: []keyboard.InlineButton{
				{
					Value: h.languageService.GetText(userLang, languageService.ContinueInterview),
					Data:  []string{command.GetNextQuestionCommand},
					Type:  keyboard.ButtonData,
				},
				{
					Value: h.languageService.GetText(userLang, languageService.StartInterviewShort),
					Data:  []string{command.ForceStartInterviewCommand},
					Type:  keyboard.ButtonData,
				},
			},
		},
	)
	if err != nil {
		return err
	}
	_, err = sender.Send(
		model.NewResponse(request.Chat.ID).
			SetText(h.languageService.GetText(userLang, languageService.ActiveInterviewExists)).
			SetInlineKeyboardMarkup(inlineKeyboard),
	)
	return err
}

func (h *Handler) Command() string {
	return command.StartInterviewCommand
}
