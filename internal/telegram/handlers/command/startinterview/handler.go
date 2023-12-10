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
	"job-interviewer/pkg/telegram/service/survey"
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

	interviewAvailableValues := h.getInterviewUC.GetAvailableValues()
	specificChainContext := chainContext{
		userLang:                 userLang,
		languageStorage:          h.languageStorage,
		interviewAvailableValues: interviewAvailableValues,
		requestData:              request.Data,
	}
	if activeInterview != nil {
		specificChainContext.activeInterviewExists = true
		specificChainContext.activeInterviewJobPosition = activeInterview.JobInfo.Position
	}

	specificChain := chain{}.
		next(activeInterviewQuestion()).
		next(mainPositionQuestion()).
		next(subPositionQuestion())

	outSurvey := out{}
	activeQuestion, err := specificChain.
		perform(
			specificChainContext,
			survey.New(),
			&outSurvey,
		).
		SetAnswers(request.Data).
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

	jobPosition := interviewAvailableValues.Nodes[outSurvey.jobPositionMainKey].Position
	if outSurvey.jobPositionSubKey != "" && len(interviewAvailableValues.Nodes[outSurvey.jobPositionMainKey].Children) > 0 {
		jobPosition = interviewAvailableValues.Nodes[outSurvey.jobPositionMainKey].Children[outSurvey.jobPositionSubKey].Position
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
