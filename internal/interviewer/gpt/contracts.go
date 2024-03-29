package gpt

import (
	"context"
	"github.com/sashabaranov/go-openai"
	"job-interviewer/internal/interviewer/model"
)

type (
	Role string

	externalClient interface {
		CreateChatCompletion(
			ctx context.Context,
			request openai.ChatCompletionRequest,
		) (response openai.ChatCompletionResponse, err error)
	}

	AcceptAnswerIn struct {
		Answer   string
		Question string
	}

	GetQuestionsListIn struct {
		JobPosition   string
		QuestionCount int64
	}

	Gateway interface {
		StartDialogue(ctx context.Context, startPrompt string) (*model.Message, error)
		ContinueDialogue(ctx context.Context, dialog []model.Message) (*model.Message, error)
		SummarizeDialogue(ctx context.Context, dialog []model.Message) (*model.Message, error)
		GetAnswerSuggestion(ctx context.Context, dialog []model.Message) (*model.Message, error)
	}
)
