package gpt

import (
	"context"
	"fmt"
	openai "github.com/sashabaranov/go-openai"
	"regexp"
	"strings"
)

const (
	getQuestionPrompt              = "I want train job interview for %s position. Ask me a question, please text question only"
	acceptAnswerPrompt             = "Evaluate my answer: \"%s\" on interviewer question and give some little tips."
	summarizeAnswersCommentsPrompt = "Make total summary all messages"
)

type DefaultGateway struct {
	client externalClient
}

func NewGateway(c externalClient) *DefaultGateway {
	return &DefaultGateway{client: c}
}

func (g *DefaultGateway) GetQuestion(ctx context.Context, jobPosition string) (string, error) {
	response, err := g.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:            openai.GPT3Dot5Turbo,
			MaxTokens:        150,
			TopP:             1,
			Temperature:      0.5,
			FrequencyPenalty: 0,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: fmt.Sprintf(getQuestionPrompt, jobPosition),
				},
			},
		},
	)
	if err != nil {
		return "", err
	}

	if len(response.Choices) == 0 {
		return "", nil
	}

	content := regexp.
		MustCompile(`[0-9]. `).
		ReplaceAllString(response.Choices[0].Message.Content, "")
	return content, nil
}

func (g *DefaultGateway) AcceptAnswer(ctx context.Context, in AcceptAnswerIn) (string, error) {
	response, err := g.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:       openai.GPT3Dot5Turbo0301,
			MaxTokens:   200,
			TopP:        0.7,
			Temperature: 0.5,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleAssistant,
					Content: fmt.Sprintf("Interviewer Question: %s", in.Question),
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: fmt.Sprintf(acceptAnswerPrompt, in.Answer),
				},
			},
		},
	)
	if err != nil {
		return "", err
	}

	out := strings.Replace(response.Choices[0].Message.Content, "\n", "", -1)
	return out, nil
}

func (g *DefaultGateway) SummarizeAnswersComments(ctx context.Context, answersComments []string) (string, error) {
	messages := make([]openai.ChatCompletionMessage, 0, len(answersComments)+1)
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: summarizeAnswersCommentsPrompt,
	})
	for _, answersComment := range answersComments {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: answersComment,
		})
	}

	response, err := g.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:       openai.GPT3Dot5Turbo0301,
			MaxTokens:   200,
			TopP:        1,
			Temperature: 0.5,
			Messages:    messages,
		},
	)
	if err != nil {
		return "", err
	}

	out := strings.Replace(response.Choices[0].Message.Content, "\n", "", -1)
	return out, nil
}
