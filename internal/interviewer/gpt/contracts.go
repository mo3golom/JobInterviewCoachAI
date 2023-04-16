package gpt

import (
	"context"
	"github.com/sashabaranov/go-openai"
)

type (
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
		GetQuestion(ctx context.Context, jobPosition string) (string, error)
		AcceptAnswer(ctx context.Context, in AcceptAnswerIn) (string, error)
		SummarizeAnswersComments(ctx context.Context, answersComments []string) (string, error)
	}
)
