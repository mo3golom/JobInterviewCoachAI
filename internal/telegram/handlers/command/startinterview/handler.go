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
	"sort"
	"strconv"
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
	var jobPositionMainKey, jobPositionSubKey string
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
					h.languageStorage.GetText(userLang, language.TextKey(strconv.FormatInt(QuestionContinueActiveInterview, 10))),
					activeInterview.JobInfo.Position,
				),
				func(answer bool) {
					skipActiveInterview = !answer
					skipActiveInterviewAnswered = true
				},
				survey.NewComplexPossibleAnswer(h.languageStorage.GetText(userLang, textKeyYes), true),
				survey.NewComplexPossibleAnswer(h.languageStorage.GetText(userLang, textKeyNo), false),
			),
		)
	}

	interviewAvailableValues := h.getInterviewUC.GetAvailableValues()
	newSurvey = newSurvey.AddQuestion(
		survey.NewQuestion(
			h.languageStorage.GetText(userLang, language.TextKey(strconv.FormatInt(QuestionJobPosition, 10))),
			func(answer string) {
				jobPositionMainKey = answer
			},
			func() []survey.PossibleAnswer[string] {
				possibleAnswers := make([]survey.PossibleAnswer[string], 0, len(interviewAvailableValues.Nodes))
				for key := range interviewAvailableValues.Nodes {
					possibleAnswers = append(
						possibleAnswers,
						survey.NewComplexPossibleAnswer(
							h.languageStorage.GetText(userLang, language.TextKey(key)),
							key,
						),
					)
				}
				sort.Slice(possibleAnswers, func(i, j int) bool {
					return possibleAnswers[i].GetContent() > possibleAnswers[j].GetContent()
				})

				return possibleAnswers
			}()...,
		),
	).SetAnswers(request.Data)

	if mainPosition, ok := interviewAvailableValues.Nodes[jobPositionMainKey]; ok && len(mainPosition.Children) > 0 {
		newSurvey = newSurvey.AddQuestion(
			survey.NewQuestion(
				fmt.Sprintf(
					h.languageStorage.GetText(userLang, textKeyClarifyJobPosition),
					h.languageStorage.GetText(userLang, language.TextKey(jobPositionMainKey)),
				),
				func(answer string) {
					jobPositionSubKey = answer
				},
				func() []survey.PossibleAnswer[string] {
					possibleAnswers := make([]survey.PossibleAnswer[string], 0, len(mainPosition.Children))
					for key := range mainPosition.Children {
						possibleAnswers = append(
							possibleAnswers,
							survey.NewComplexPossibleAnswer(
								h.languageStorage.GetText(userLang, language.TextKey(key)),
								key,
							),
						)
					}
					sort.Slice(possibleAnswers, func(i, j int) bool {
						return possibleAnswers[i].GetContent() > possibleAnswers[j].GetContent()
					})

					return possibleAnswers
				}()...,
			),
		)
	}

	activeQuestion, err := newSurvey.
		SetAnswers(request.Data).
		FindUnansweredQuestionAsKeyboard(h.Command())
	if err != nil {
		return err
	}

	if skipActiveInterviewAnswered && !skipActiveInterview {
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

	jobPosition := interviewAvailableValues.Nodes[jobPositionMainKey].Position
	if jobPositionSubKey != "" && len(interviewAvailableValues.Nodes[jobPositionMainKey].Children) > 0 {
		jobPosition = interviewAvailableValues.Nodes[jobPositionMainKey].Children[jobPositionSubKey].Position
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
