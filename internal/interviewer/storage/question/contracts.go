package question

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"job-interviewer/internal/interviewer/model"
	"job-interviewer/pkg/transactional"
)

var (
	ErrEmptyQuestionResult = errors.New("empty question result")
)

type (
	SetQuestionAnsweredIn struct {
		InterviewID uuid.UUID
		QuestionID  uuid.UUID
		Answer      string
		GptComment  string
	}

	UpdateInterviewQuestionStatusIn struct {
		InterviewID uuid.UUID
		Current     model.InterviewQuestionStatus
		Target      model.InterviewQuestionStatus
	}

	Storage interface {
		CreateQuestions(ctx context.Context, tx transactional.Tx, in []model.Question) error
		AttachQuestionsToInterview(ctx context.Context, tx transactional.Tx, interviewID uuid.UUID, questions []model.Question) error
		SetQuestionAnswered(ctx context.Context, tx transactional.Tx, in SetQuestionAnsweredIn) error
		UpdateInterviewQuestionStatus(ctx context.Context, tx transactional.Tx, in UpdateInterviewQuestionStatusIn) error

		FindNextQuestion(ctx context.Context, tx transactional.Tx, interviewID uuid.UUID) (*model.Question, error)
		FindActiveQuestionByInterviewID(ctx context.Context, interviewID uuid.UUID) (*model.Question, error)
	}
)
