package interview

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"job-interviewer/internal/interviewer/model"
)

var (
	ErrAlreadyExistsStartedInterview = errors.New("already exists started interview")
)

type (
	CreateInterviewIn struct {
		UserID         uuid.UUID
		JobPosition    string
		QuestionsCount int64
	}

	AcceptAnswerIn struct {
		UserID uuid.UUID
		Answer string
	}

	Service interface {
		CreateInterview(ctx context.Context, in CreateInterviewIn) (*model.Interview, error)
		StartInterview(ctx context.Context, interview *model.Interview) error
		FinishInterview(ctx context.Context, interview *model.Interview) (string, error)
		FinishInterviewWithoutSummary(ctx context.Context, interview *model.Interview) error
		FindActiveInterview(ctx context.Context, userID uuid.UUID) (*model.Interview, error)
		AcceptAnswer(ctx context.Context, in AcceptAnswerIn) (string, error)
	}
)
