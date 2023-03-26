package startinterview

import (
	"context"
	"fmt"
	interviewContracts "job-interviewer/internal/interviewer/contracts"
	"job-interviewer/internal/telegram/handlers/command"
	languageService "job-interviewer/internal/telegram/language"
	"job-interviewer/internal/telegram/service"
	"job-interviewer/pkg/telegram"
	"job-interviewer/pkg/telegram/model"
	"job-interviewer/pkg/telegram/service/keyboard"
	"sort"
)

type Handler struct {
	keyboardService       keyboard.Service
	getInterviewOptionsUC interviewContracts.GetInterviewOptionsUseCase
	startInterviewUC      interviewContracts.StartInterviewUseCase
	service               service.Service
	languageService       languageService.Service
}

func NewHandler(
	k keyboard.Service,
	g interviewContracts.GetInterviewOptionsUseCase,
	s interviewContracts.StartInterviewUseCase,
	service service.Service,
	l languageService.Service,
) *Handler {
	return &Handler{
		keyboardService:       k,
		getInterviewOptionsUC: g,
		startInterviewUC:      s,
		service:               service,
		languageService:       l,
	}
}

func (h *Handler) Handle(ctx context.Context, request *model.Request, sender telegram.Sender) error {
	if len(request.Data) == 0 {
		return h.choosePosition(request, sender)
	}

	if len(request.Data) == 1 {
		return h.chooseLevel(request, sender)

	}

	return h.startInterview(ctx, request, sender)
}

func (h *Handler) choosePosition(request *model.Request, sender telegram.Sender) error {
	userLang := request.User.Lang

	interviewOptions := h.getInterviewOptionsUC.GetInterviewOptions()
	buttons := make([]keyboard.InlineButton, 0, len(interviewOptions.Positions))
	for key, position := range interviewOptions.Positions {
		buttons = append(
			buttons,
			keyboard.InlineButton{
				Value: position,
				Data:  []string{key},
				Type:  keyboard.ButtonData,
			},
		)
	}

	sort.Slice(buttons, func(i, j int) bool {
		return buttons[i].Value < buttons[j].Value
	})

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
			SetText(h.languageService.GetText(userLang, languageService.ChoosePosition)).
			SetInlineKeyboardMarkup(inlineKeyboard),
	)
	return err
}

func (h *Handler) chooseLevel(request *model.Request, sender telegram.Sender) error {
	userLang := request.User.Lang
	position := request.Data[0]

	currentCommand := h.Command()
	interviewOptions := h.getInterviewOptionsUC.GetInterviewOptions()
	buttons := make([]keyboard.InlineButton, 0, len(interviewOptions.Levels))
	for key, level := range interviewOptions.Levels {
		buttons = append(
			buttons,
			keyboard.InlineButton{
				Value: levelToString[level],
				Data:  []string{position, key},
				Type:  keyboard.ButtonData,
			},
		)
	}

	sort.Slice(buttons, func(i, j int) bool {
		return buttons[i].Value < buttons[j].Value
	})

	inlineKeyboard, err := h.keyboardService.BuildInlineKeyboardGrid(
		keyboard.BuildInlineKeyboardIn{
			Command: &currentCommand,
			Buttons: buttons,
		},
	)
	if err != nil {
		return err
	}
	return sender.Update(
		request.MessageID,
		model.NewResponse(request.Chat.ID).
			SetText(h.languageService.GetText(userLang, languageService.ChooseLevel)).
			SetInlineKeyboardMarkup(inlineKeyboard),
	)
}

func (h *Handler) startInterview(ctx context.Context, request *model.Request, sender telegram.Sender) error {
	userLang := request.User.Lang
	interviewOptions := h.getInterviewOptionsUC.GetInterviewOptions()
	position := interviewOptions.Positions[request.Data[0]]
	level := interviewOptions.Levels[request.Data[1]]

	err := sender.Update(
		request.MessageID,
		model.NewResponse(request.Chat.ID).
			SetText(
				fmt.Sprintf(
					h.languageService.GetText(userLang, languageService.StartInterviewSummary),
					position,
					levelToString[level],
				),
			).
			SetInlineKeyboardMarkup(nil),
	)
	if err != nil {
		return err
	}

	_, err = sender.Send(model.NewResponse(request.Chat.ID).SetText(
		h.languageService.GetText(userLang, languageService.LoadQuestions)),
	)
	if err != nil {
		return err
	}

	err = h.startInterviewUC.StartInterview(
		ctx,
		interviewContracts.StartInterviewIn{
			UserID:         request.User.OriginalID,
			JobPosition:    position,
			JobLevel:       level,
			QuestionsCount: 10, //TODO: добавить выбор количества вопросов
		},
	)
	if err != nil {
		return err
	}

	return h.service.GetNextQuestion(ctx, request, sender)
}

func (h *Handler) Command() string {
	return command.ForceStartInterviewCommand
}
