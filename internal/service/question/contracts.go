package question

import (
	"context"
	"github.com/google/uuid"
	"job-interviewer/internal/model"
)

type Service interface {
	CreateQuestionsForInterview(ctx context.Context, interview *model.Interview) error
	FindNextQuestion(ctx context.Context, interviewID uuid.UUID) (*model.Question, error)
}
