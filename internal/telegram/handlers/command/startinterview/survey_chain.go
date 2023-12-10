package startinterview

import (
	"fmt"
	"job-interviewer/internal/interviewer/model"
	"job-interviewer/pkg/language"
	"job-interviewer/pkg/telegram/service/survey"
	"sort"
)

type (
	chainContext struct {
		activeInterviewExists      bool
		activeInterviewJobPosition model.Position
		userLang                   language.Language
		languageStorage            language.Storage
		requestData                []string
		interviewAvailableValues   *model.InterviewAvailableValues
	}

	out struct {
		skipActiveInterview         bool
		skipActiveInterviewAnswered bool
		jobPositionMainKey          string
		jobPositionSubKey           string
	}

	chain struct {
		items []chainItem
	}

	chainItem struct {
		fn func(ctx chainContext, mainSurvey survey.Survey, out *out) survey.Survey
	}
)

func (c chain) next(nextChainItem chainItem) chain {
	c.items = append(c.items, nextChainItem)

	return c
}

func (c chain) perform(ctx chainContext, mainSurvey survey.Survey, out *out) survey.Survey {
	for _, item := range c.items {
		mainSurvey = item.fn(ctx, mainSurvey, out)
	}

	return mainSurvey
}

func activeInterviewQuestion() chainItem {
	return chainItem{
		fn: func(ctx chainContext, mainSurvey survey.Survey, out *out) survey.Survey {
			if !ctx.activeInterviewExists {
				return mainSurvey
			}

			return mainSurvey.AddQuestion(
				survey.NewQuestion(
					fmt.Sprintf(
						ctx.languageStorage.GetText(ctx.userLang, textKeyQuestionContinueActiveInterview),
						ctx.activeInterviewJobPosition,
					),
					func(answer bool) {
						out.skipActiveInterview = !answer
						out.skipActiveInterviewAnswered = true
					},
					survey.NewComplexPossibleAnswer(ctx.languageStorage.GetText(ctx.userLang, textKeyYes), true),
					survey.NewComplexPossibleAnswer(ctx.languageStorage.GetText(ctx.userLang, textKeyNo), false),
				),
			)
		},
	}
}

func mainPositionQuestion() chainItem {
	return chainItem{
		fn: func(ctx chainContext, mainSurvey survey.Survey, out *out) survey.Survey {
			return mainSurvey.AddQuestion(
				survey.NewQuestion(
					ctx.languageStorage.GetText(ctx.userLang, textKeyQuestionJobPosition),
					func(answer string) {
						out.jobPositionMainKey = answer
					},
					func() []survey.PossibleAnswer[string] {
						possibleAnswers := make([]survey.PossibleAnswer[string], 0, len(ctx.interviewAvailableValues.Nodes))
						for key := range ctx.interviewAvailableValues.Nodes {
							possibleAnswers = append(
								possibleAnswers,
								survey.NewComplexPossibleAnswer(
									ctx.languageStorage.GetText(ctx.userLang, language.TextKey(key)),
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
			).SetAnswers(ctx.requestData)
		},
	}
}

func subPositionQuestion() chainItem {
	return chainItem{
		fn: func(ctx chainContext, mainSurvey survey.Survey, out *out) survey.Survey {
			mainPosition, ok := ctx.interviewAvailableValues.Nodes[out.jobPositionMainKey]
			if !ok || len(mainPosition.Children) <= 0 {
				return mainSurvey
			}

			return mainSurvey.AddQuestion(
				survey.NewQuestion(
					fmt.Sprintf(
						ctx.languageStorage.GetText(ctx.userLang, textKeyClarifyJobPosition),
						ctx.languageStorage.GetText(ctx.userLang, language.TextKey(out.jobPositionMainKey)),
					),
					func(answer string) {
						out.jobPositionSubKey = answer
					},
					func() []survey.PossibleAnswer[string] {
						possibleAnswers := make([]survey.PossibleAnswer[string], 0, len(mainPosition.Children))
						for key := range mainPosition.Children {
							possibleAnswers = append(
								possibleAnswers,
								survey.NewComplexPossibleAnswer(
									ctx.languageStorage.GetText(ctx.userLang, language.TextKey(key)),
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
		},
	}
}
