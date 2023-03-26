package changeuserlanguage

import (
	"context"
	"job-interviewer/internal/interviewer/contracts"
	"job-interviewer/internal/telegram/handlers/command"
	languageService "job-interviewer/internal/telegram/language"
	"job-interviewer/internal/telegram/service"
	"job-interviewer/pkg/language"
	"job-interviewer/pkg/telegram"
	"job-interviewer/pkg/telegram/model"
	"job-interviewer/pkg/telegram/service/keyboard"
)

type Handler struct {
	userUC          contracts.UserUseCase
	service         service.Service
	keyboardService keyboard.Service
	languageService languageService.Service
}

func NewHandler(
	userUC contracts.UserUseCase,
	service service.Service,
	keyboardService keyboard.Service,
	languageService languageService.Service,
) *Handler {
	return &Handler{
		userUC:          userUC,
		service:         service,
		keyboardService: keyboardService,
		languageService: languageService,
	}
}

func (h *Handler) Handle(ctx context.Context, request *model.Request, sender telegram.Sender) error {
	if len(request.Data) == 0 {
		return h.chooseLanguage(request, sender)
	}

	userLang := language.Language(request.Data[0])
	err := h.userUC.ChangeLanguage(
		ctx,
		request.User.OriginalID,
		userLang,
	)
	if err != nil {
		return err
	}

	_, err = sender.Send(
		model.NewResponse(request.Chat.ID).
			SetText(h.languageService.GetText(userLang, languageService.ChooseLanguageSuccess)).
			SetKeyboardMarkup(h.service.GetUserMainKeyboard(userLang)),
	)
	return err
}

func (h *Handler) Command() string {
	return command.ChangeUserLanguage
}

func (h *Handler) chooseLanguage(request *model.Request, sender telegram.Sender) error {
	userLang := request.User.Lang

	currentCommand := h.Command()
	inlineKeyboard, err := h.keyboardService.BuildInlineKeyboardGrid(
		keyboard.BuildInlineKeyboardIn{
			Command: &currentCommand,
			Buttons: buttons,
		},
	)
	if err != nil {
		return err
	}

	_, err = sender.Send(
		model.NewResponse(request.Chat.ID).
			SetText(h.languageService.GetText(userLang, languageService.ChooseLanguage)).
			SetInlineKeyboardMarkup(inlineKeyboard),
	)
	return err
}
