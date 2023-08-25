package gpt

import (
	"context"
	"fmt"
	openai "github.com/sashabaranov/go-openai"
	"job-interviewer/internal/interviewer/model"
	"regexp"
	"strings"
)

const (
	getQuestionPrompt              = `I want you to act as an interviewer. I will be the candidate and you will ask me the interview questions for the %s position. I want you to only reply as the interviewer. Do not write all the conservation at once. I want you to only do the interview with me. Ask me the tricky questions and wait for my answers. Do not write explanations. Do not write "interviewer:". Ask me the questions one by one like an interviewer does and wait for my answers. My first sentence is "Hi"`
	getPossibleAnswersPrompt       = "I don't know how to answer to this question, please give me a list of possible answers"
	summarizeAnswersCommentsPrompt = "I want to finish interview. Summarize the dialogue and give me feedback. Ignore messages where I ask help with answer."
)

type DefaultGateway struct {
	client externalClient
}

func NewGateway(c externalClient) *DefaultGateway {
	return &DefaultGateway{client: c}
}

func (g *DefaultGateway) SummarizeAnswersComments(ctx context.Context, dialog []model.Message, jobPosition string) (*model.Message, error) {
	messages := make([]openai.ChatCompletionMessage, 0, len(dialog))
	for _, message := range dialog {
		role := openai.ChatMessageRoleAssistant
		if message.Role == model.RoleUser {
			role = openai.ChatMessageRoleUser
		}

		messages = append(
			messages,
			openai.ChatCompletionMessage{
				Role:    role,
				Content: message.Content,
			},
		)
	}

	messages = append(
		messages,
		openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: summarizeAnswersCommentsPrompt,
		},
	)

	return g.createChatCompletion(
		ctx,
		append([]openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: fmt.Sprintf(getQuestionPrompt, jobPosition),
			},
		}, messages...),
	)
}

func (g *DefaultGateway) StartDialogue(ctx context.Context, jobPosition string) (*model.Message, error) {
	return g.createChatCompletion(
		ctx,
		[]openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: fmt.Sprintf(getQuestionPrompt, jobPosition),
			},
		},
	)
}

func (g *DefaultGateway) ContinueDialogue(ctx context.Context, dialog []model.Message, jobPosition string) (*model.Message, error) {
	messages := make([]openai.ChatCompletionMessage, 0, len(dialog))
	for _, message := range dialog {
		role := openai.ChatMessageRoleAssistant
		if message.Role == model.RoleUser {
			role = openai.ChatMessageRoleUser
		}

		messages = append(
			messages,
			openai.ChatCompletionMessage{
				Role:    role,
				Content: message.Content,
			},
		)
	}

	return g.createChatCompletion(
		ctx,
		append([]openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: fmt.Sprintf(getQuestionPrompt, jobPosition),
			},
		}, messages...),
	)
}

func (g *DefaultGateway) createChatCompletion(ctx context.Context, messages []openai.ChatCompletionMessage) (*model.Message, error) {
	response, err := g.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:       openai.GPT3Dot5Turbo,
			MaxTokens:   200,
			TopP:        0.7,
			Temperature: 0.5,
			Messages:    messages,
		},
	)
	if err != nil {
		return nil, err
	}

	if len(response.Choices) == 0 {
		return nil, nil
	}

	content := regexp.
		MustCompile(`[0-9]. `).
		ReplaceAllString(response.Choices[0].Message.Content, "")
	content = strings.TrimSpace(content)
	return &model.Message{
		Role:    model.RoleAssistant,
		Content: content,
	}, err
}
