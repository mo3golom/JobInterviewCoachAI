package startinterview

import (
	"context"
	"fmt"
	"job-interviewer/internal/contracts"
	modelInterview "job-interviewer/internal/model"
	"job-interviewer/pkg/telegram"
	"job-interviewer/pkg/telegram/handlers/command"
	"job-interviewer/pkg/telegram/model"
	"job-interviewer/pkg/telegram/service/keyboard"
)

type Handler struct {
	keyboardService        keyboard.Service
	getInterviewOptionsUC  contracts.GetInterviewOptionsUseCase
	startInterviewUC       contracts.StartInterviewUseCase
	getNextQuestionHandler telegram.Handler
}

func NewHandler(
	k keyboard.Service,
	g contracts.GetInterviewOptionsUseCase,
	s contracts.StartInterviewUseCase,
	pr telegram.Handler,
) *Handler {
	return &Handler{
		keyboardService:        k,
		getInterviewOptionsUC:  g,
		startInterviewUC:       s,
		getNextQuestionHandler: pr,
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
	interviewOptions := h.getInterviewOptionsUC.GetInterviewOptions()
	buttons := make([]keyboard.InlineButton, 0, len(interviewOptions.Positions))
	for _, position := range interviewOptions.Positions {
		buttons = append(
			buttons,
			keyboard.InlineButton{
				Value: position,
				Data:  []string{position},
				Type:  keyboard.ButtonData,
			},
		)
	}

	currentCommand := h.Command()
	_, err := sender.Send(
		model.NewResponse(request.Chat.ID).
			SetText(choosePositionText).
			SetInlineKeyboardMarkup(
				h.keyboardService.BuildInlineKeyboardGrid(
					keyboard.BuildInlineKeyboardIn{
						Command: &currentCommand,
						Buttons: buttons,
					},
				),
			),
	)
	return err
}

func (h *Handler) chooseLevel(request *model.Request, sender telegram.Sender) error {
	position := request.Data[0]

	currentCommand := h.Command()
	interviewOptions := h.getInterviewOptionsUC.GetInterviewOptions()
	buttons := make([]keyboard.InlineButton, 0, len(interviewOptions.Levels))
	for _, level := range interviewOptions.Levels {
		buttons = append(
			buttons,
			keyboard.InlineButton{
				Value: levelToString[level],
				Data:  []string{position, string(level)},
				Type:  keyboard.ButtonData,
			},
		)
	}

	return sender.Update(
		request.MessageID,
		model.NewResponse(request.Chat.ID).
			SetText(chooseLevelText).
			SetInlineKeyboardMarkup(
				h.keyboardService.BuildInlineKeyboardGrid(
					keyboard.BuildInlineKeyboardIn{
						Command: &currentCommand,
						Buttons: buttons,
					},
				),
			),
	)
}

func (h *Handler) startInterview(ctx context.Context, request *model.Request, sender telegram.Sender) error {
	position := request.Data[0]
	level := modelInterview.JobLevel(request.Data[1])

	err := sender.Update(
		request.MessageID,
		model.NewResponse(request.Chat.ID).
			SetText(
				fmt.Sprintf(
					startInterviewText,
					position,
					levelToString[level],
				),
			).
			SetInlineKeyboardMarkup(nil),
	)
	if err != nil {
		return err
	}

	_, err = sender.Send(model.NewResponse(request.Chat.ID).SetText(loadQuestionsText))
	if err != nil {
		return err
	}

	err = h.startInterviewUC.StartInterview(
		ctx,
		contracts.StartInterviewIn{
			UserID:         request.User.OriginalID,
			JobPosition:    position,
			JobLevel:       level,
			QuestionsCount: 10, //TODO: добавить выбор количества вопросов
		},
	)
	if err != nil {
		return err
	}

	return h.getNextQuestionHandler.Handle(ctx, request, sender)
}

func (h *Handler) Command() string {
	return command.ForceStartInterviewCommand
}
