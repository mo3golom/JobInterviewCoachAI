package gpt

import (
	"context"
	openai "github.com/sashabaranov/go-openai"
	"job-interviewer/internal/interviewer/model"
	"regexp"
	"strings"
)

const (
	getAnswerSuggestionPrompt      = "I don't know how to answer to this question, please give me a list of possible answers"
	summarizeAnswersCommentsPrompt = "I want to finish interview. Summarize the dialogue and give me feedback. Ignore messages where I ask help with answer."
)

type DefaultGateway struct {
	client externalClient
}

func NewGateway(c externalClient) *DefaultGateway {
	return &DefaultGateway{client: c}
}

func (g *DefaultGateway) SummarizeDialogue(ctx context.Context, dialog []model.Message) (*model.Message, error) {
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
		messages,
	)
}

func (g *DefaultGateway) StartDialogue(ctx context.Context, startPrompt string) (*model.Message, error) {
	return g.createChatCompletion(
		ctx,
		[]openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: startPrompt,
			},
		},
	)
}

func (g *DefaultGateway) ContinueDialogue(ctx context.Context, dialog []model.Message) (*model.Message, error) {
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
		messages,
	)
}

func (g *DefaultGateway) GetAnswerSuggestion(ctx context.Context, dialog []model.Message) (*model.Message, error) {
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
			Content: getAnswerSuggestionPrompt,
		},
	)

	return g.createChatCompletion(
		ctx,
		messages,
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
	content = strings.TrimSuffix(strings.TrimSpace(content), "\n")
	return &model.Message{
		Role:    model.RoleAssistant,
		Content: content,
	}, err
}
