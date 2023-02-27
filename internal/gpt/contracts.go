package gpt

import (
	"context"
	gogpt "github.com/sashabaranov/go-gpt3"
)

type externalClient interface {
	CreateCompletion(ctx context.Context, request gogpt.CompletionRequest) (response gogpt.CompletionResponse, err error)
	CreateCompletionStream(ctx context.Context, request gogpt.CompletionRequest) (stream *gogpt.CompletionStream, err error)
}

type AcceptAnswerIn struct {
	Answer      string
	Question    string
	JobPosition string
}

type GetQuestionsListIn struct {
	JobPosition   string
	QuestionCount int64
}

type Gateway interface {
	AcceptAnswer(ctx context.Context, in AcceptAnswerIn) (string, error)
	GetQuestionsList(ctx context.Context, in GetQuestionsListIn) ([]string, error)
}
