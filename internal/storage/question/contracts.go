package question

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"job-interviewer/internal/model"
	"job-interviewer/pkg/transactional"
)

var (
	ErrEmptyQuestionResult = errors.New("empty question result")
)

type Storage interface {
	CreateQuestions(ctx context.Context, tx transactional.Tx, in []model.Question) error
	AttachQuestionsToInterview(ctx context.Context, tx transactional.Tx, interviewID uuid.UUID, questions []model.Question) error
	SetQuestionAnswered(ctx context.Context, tx transactional.Tx, interviewID uuid.UUID, questionID uuid.UUID) error

	FindNextQuestion(ctx context.Context, tx transactional.Tx, interviewID uuid.UUID) (*model.Question, error)
	FindActiveQuestionByInterviewID(ctx context.Context, interviewID uuid.UUID) (*model.Question, error)
}
