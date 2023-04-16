package question

import (
	"context"
	"job-interviewer/internal/interviewer/model"
)

type Service interface {
	GetNextQuestion(ctx context.Context, interview *model.Interview) (*model.Question, error)
}
