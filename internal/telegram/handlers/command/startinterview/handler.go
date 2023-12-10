package startinterview

import (
	"context"
	"errors"
	"fmt"
	interviewContracts "job-interviewer/internal/interviewer/contracts"
	"job-interviewer/internal/telegram/handlers/command"
	"job-interviewer/internal/telegram/service"
	"job-interviewer/pkg/language"
	"job-interviewer/pkg/telegram"
	"job-interviewer/pkg/telegram/model"
)

type Handler struct {
	startInterviewUC interviewContracts.StartInterviewUseCase
	getInterviewUC   interviewContracts.GetInterviewUsecase
	service          service.Service
	languageStorage  language.Storage
}

func NewHandler(
	s interviewContracts.StartInterviewUseCase,
	getInterviewUC interviewContracts.GetInterviewUsecase,
	service service.Service,
) *Handler {
	return &Handler{
		startInterviewUC: s,
		getInterviewUC:   getInterviewUC,
		service:          service,
		languageStorage:  configLanguage(),
	}
}

func (h *Handler) Handle(ctx context.Context, request *model.Request, sender telegram.Sender) error {
	userLang := language.Russian

	activeInterview, err := h.getInterviewUC.FindActiveInterview(ctx, request.User.OriginalID)
	if err != nil && !errors.Is(err, interviewContracts.ErrEmptyActiveInterview) {
		return err
	}

	outSurvey := out{}
	activeQuestion, err := h.getSurvey(
		userLang,
		activeInterview,
		request.Data,
		&outSurvey,
	).
		FindUnansweredQuestionAsKeyboard(h.Command())
	if err != nil {
		return err
	}

	if outSurvey.skipActiveInterviewAnswered && !outSurvey.skipActiveInterview {
		err = h.startInterviewUC.ContinueInterview(ctx, request.User.OriginalID)
		if err != nil {
			return err
		}

		return h.service.GetNextQuestion(ctx, request, sender)
	}

	if activeQuestion != nil {
		_, err = sender.Send(
			model.NewResponse().
				SetText(activeQuestion.Text).
				SetInlineKeyboardMarkup(activeQuestion.Keyboard),
		)
		return err
	}

	jobPosition := outSurvey.jobPositionMainKey
	if outSurvey.jobPositionSubKey != "" {
		jobPosition = outSurvey.jobPositionSubKey
	}

	err = h.startInterviewUC.StartInterview(
		ctx,
		interviewContracts.StartInterviewIn{
			UserID: request.User.OriginalID,
			Questions: interviewContracts.StartInterviewQuestionsIn{
				JobPosition: jobPosition,
			},
		},
	)
	if err != nil {
		return err
	}

	finalJobPositionKey := outSurvey.jobPositionMainKey
	if outSurvey.jobPositionSubKey != "" {
		finalJobPositionKey = outSurvey.jobPositionSubKey
	}
	_, err = sender.Send(
		model.NewResponse().
			SetText(fmt.Sprintf(
				h.languageStorage.GetText(userLang, textKeyYourChoice),
				h.languageStorage.GetText(userLang, language.TextKey(finalJobPositionKey)),
			)),
	)
	if err != nil {
		return err
	}

	return h.service.GetNextQuestion(ctx, request, sender)
}

func (h *Handler) Command() string {
	return command.ForceStartInterviewCommand
}

func (h *Handler) Aliases() []string {
	return []string{
		h.languageStorage.GetText(language.Russian, textKeyStartInterview),
	}
}
