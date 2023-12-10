package startinterview

import (
	"fmt"
	"job-interviewer/internal"
	"job-interviewer/internal/interviewer/model"
	"job-interviewer/pkg/language"
	"job-interviewer/pkg/telegram/survey"
	"sort"
)

var (
	MainPositions = []internal.Position{
		Developer,
		internal.ProjectManager,
		internal.ProductManager,
		internal.ProductDesigner,
		internal.Behavioral,
	}

	SubPositions = map[internal.Position][]internal.Position{
		Developer: {
			internal.GolangDeveloper,
			internal.PhpDeveloper,
			internal.PythonDeveloper,
			internal.RustDeveloper,
			internal.JavascriptDeveloper,
			internal.SwiftDeveloper,
			internal.JavaDeveloper,
			internal.CplusplusDeveloper,
			internal.CsharpDeveloper,
		},
	}
)

type (
	out struct {
		skipActiveInterview         bool
		skipActiveInterviewAnswered bool
		jobPositionMainKey          internal.Position
		jobPositionSubKey           internal.Position
	}
)

func (h *Handler) getSurvey(
	userLang language.Language,
	activeInterview *model.Interview,
	requestData []string,
	outSurvey *out,
) survey.Survey[*out] {
	var activeInterviewPosition internal.Position
	if activeInterview != nil {
		activeInterviewPosition = activeInterview.JobInfo.Position
	}

	return survey.New[*out]().
		AddQuestionWhen(
			survey.NewQuestion(
				fmt.Sprintf(
					h.languageStorage.GetText(userLang, textKeyQuestionContinueActiveInterview),
					h.languageStorage.GetText(userLang, language.TextKey(activeInterviewPosition)),
				),
				func(answer bool, out *out) {
					out.skipActiveInterview = !answer
					out.skipActiveInterviewAnswered = true
				},
				survey.NewComplexPossibleAnswer(h.languageStorage.GetText(userLang, textKeyYes), true),
				survey.NewComplexPossibleAnswer(h.languageStorage.GetText(userLang, textKeyNo), false),
			),
			func() bool {
				return activeInterview != nil
			},
		).AddQuestion(
		survey.NewQuestion(
			h.languageStorage.GetText(userLang, textKeyQuestionJobPosition),
			func(answer internal.Position, out *out) {
				out.jobPositionMainKey = answer
			},
			func() []survey.PossibleAnswer[internal.Position] {
				possibleAnswers := make([]survey.PossibleAnswer[internal.Position], 0, len(MainPositions))
				for _, key := range MainPositions {
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
		)).
		Init(requestData, outSurvey).
		AddQuestionWhen(
			survey.NewQuestion(
				fmt.Sprintf(
					h.languageStorage.GetText(userLang, textKeyClarifyJobPosition),
					h.languageStorage.GetText(userLang, language.TextKey(outSurvey.jobPositionMainKey)),
				),
				func(answer internal.Position, out *out) {
					out.jobPositionSubKey = answer
				},
				func() []survey.PossibleAnswer[internal.Position] {
					possibleAnswers := make([]survey.PossibleAnswer[internal.Position], 0, len(SubPositions[outSurvey.jobPositionMainKey]))
					for _, key := range SubPositions[outSurvey.jobPositionMainKey] {
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
			func() bool {
				_, ok := SubPositions[outSurvey.jobPositionMainKey]

				return ok
			},
		).
		Init(requestData, outSurvey)
}
