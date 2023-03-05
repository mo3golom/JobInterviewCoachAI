package question

import (
	"context"
	"github.com/google/uuid"
	model2 "job-interviewer/internal/interviewer/model"
)

type Service interface {
	CreateQuestionsForInterview(ctx context.Context, interview *model2.Interview) error
	FindNextQuestion(ctx context.Context, interviewID uuid.UUID) (*model2.Question, error)
}
