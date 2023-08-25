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
	var jobPosition string
	var skipActiveInterview, skipActiveInterviewAnswered bool
	newSurvey := survey.New()

	activeInterview, err := h.getInterviewUC.FindActiveInterview(ctx, request.User.OriginalID)
	if err != nil && !errors.Is(err, interviewContracts.ErrEmptyActiveInterview) {
		return err
	}

	if activeInterview != nil {
		newSurvey = newSurvey.AddQuestion(
			survey.NewQuestion(
				fmt.Sprintf(
					h.languageStorage.GetText(userLang, language.TextKey(QuestionContinueActiveInterview)),
					activeInterview.JobInfo.Position,
				),
				func(answer bool) {
					skipActiveInterview = !answer
					skipActiveInterviewAnswered = true
				},
				survey.NewComplexPossibleAnswer("Да", true),
				survey.NewComplexPossibleAnswer("Нет", false),
			),
		)
	}

	activeQuestion, err := newSurvey.AddQuestion(
		survey.NewQuestion(
			h.languageStorage.GetText(userLang, language.TextKey(QuestionJobPosition)),
			func(answer string) {
				jobPosition = answer
			},
			survey.NewPossibleAnswer("golang developer"),
			survey.NewPossibleAnswer("python developer"),
			survey.NewPossibleAnswer("php developer"),
		),
	).
		SetAnswers(request.Data).
		FindUnansweredQuestionAsKeyboard(h.Command())
	if err != nil {
		return err
	}

	if skipActiveInterviewAnswered && !skipActiveInterview {
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
