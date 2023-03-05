package gpt

import (
	"context"
	"fmt"
	gogpt "github.com/sashabaranov/go-gpt3"
	"regexp"
	"strings"
)

const (
	getQuestionListPrompt = "Create a list of %d questions for my job interview on %s position:"
	acceptAnswerPrompt    = "i`m trying to job interview on %s position, please evaluate my answer: \"%s\" on question: \"%s\" and give some tips."
)

type DefaultGateway struct {
	client externalClient
}

func NewGateway(c externalClient) *DefaultGateway {
	return &DefaultGateway{client: c}
}

func (g *DefaultGateway) GetQuestionsList(ctx context.Context, in GetQuestionsListIn) ([]string, error) {
	prompt := fmt.Sprintf(
		getQuestionListPrompt,
		in.QuestionCount,
		in.JobPosition,
	)

	response, err := g.client.CreateCompletion(
		ctx,
		gogpt.CompletionRequest{
			Model:            gogpt.GPT3TextDavinci003,
			MaxTokens:        150,
			TopP:             1,
			Temperature:      0.5,
			FrequencyPenalty: 0,
			Prompt:           prompt,
		},
	)
	if err != nil {
		return nil, err
	}

	if len(response.Choices) == 0 {
		return nil, nil
	}

	rawText := regexp.
		MustCompile(`[0-9]. `).
		ReplaceAllString(response.Choices[0].Text, "")
	rawQuestions := strings.Split(rawText, "\n")
	questions := make([]string, 0, len(rawQuestions))
	for _, rawQuestion := range rawQuestions {
		if rawQuestion == "" || rawQuestion == " " {
			continue
		}

		questions = append(questions, rawQuestion)
	}

	return questions, nil
}

func (g *DefaultGateway) AcceptAnswer(ctx context.Context, in AcceptAnswerIn) (string, error) {
	prompt := fmt.Sprintf(
		acceptAnswerPrompt,
		in.JobPosition,
		in.Answer,
		in.Question,
	)
	response, err := g.client.CreateCompletion(
		ctx,
		gogpt.CompletionRequest{
			Model:       gogpt.GPT3TextDavinci003,
			MaxTokens:   50,
			TopP:        1,
			Temperature: 0.7,
			Prompt:      prompt,
		},
	)
	if err != nil {
		return "", err
	}

	out := strings.Replace(response.Choices[0].Text, "\n", "", -1)
	return out, nil
}
