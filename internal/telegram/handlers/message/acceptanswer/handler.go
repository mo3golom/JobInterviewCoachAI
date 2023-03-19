package acceptanswer

import (
	"context"
	"errors"
	"fmt"
	interviewerContracts "job-interviewer/internal/interviewer/contracts"
	"job-interviewer/internal/telegram/handlers"
	"job-interviewer/internal/telegram/handlers/command"
	languageService "job-interviewer/internal/telegram/language"
	"job-interviewer/internal/telegram/service"
	"job-interviewer/pkg/language"
	"job-interviewer/pkg/telegram"
	"job-interviewer/pkg/telegram/model"
	"job-interviewer/pkg/telegram/service/keyboard"
)

type Handler struct {
	acceptAnswerUC  interviewerContracts.AcceptAnswerUseCase
	keyboardService keyboard.Service
	service         service.Service
	languageService languageService.Service
}

func NewHandler(
	auc interviewerContracts.AcceptAnswerUseCase,
	s service.Service,
	ks keyboard.Service,
	l languageService.Service,
) *Handler {
	return &Handler{
		acceptAnswerUC:  auc,
		service:         s,
		keyboardService: ks,
		languageService: l,
	}
}

func (h *Handler) Handle(ctx context.Context, request *model.Request, sender telegram.Sender) error {
	userLang := language.Language(request.User.Lang)

	if request.Message == nil {
		return nil
	}

	answerMessageID, err := sender.Send(
		model.NewResponse(request.Chat.ID).SetText(
			fmt.Sprintf(
				"%s %s",
				handlers.RobotPrefix,
				h.languageService.GetText(languageService.English, languageService.ProcessingAnswer),
			),
		),
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
		return h.service.Start(request, sender)
	}
	if err != nil {
		return err
	}

	inlineKeyboard, err := h.keyboardService.BuildInlineKeyboardInlineList(keyboard.BuildInlineKeyboardIn{
		Buttons: []keyboard.InlineButton{
			{
				Value: h.languageService.GetText(userLang, languageService.FinishInterview),
				Data:  []string{command.FinishInterviewCommand},
				Type:  keyboard.ButtonData,
			},
			{
				Value: h.languageService.GetText(userLang, languageService.ContinueInterview),
				Data:  []string{command.GetNextQuestionCommand},
				Type:  keyboard.ButtonData,
			},
		},
	})
	if err != nil {
		return err
	}
	err = sender.Update(
		answerMessageID,
		response.
			SetText(
				fmt.Sprintf("%s %s", handlers.RobotPrefix, out),
			).
			SetInlineKeyboardMarkup(inlineKeyboard),
	)

	return err
}
