package interview

import (
	"context"
	"errors"
	"github.com/google/uuid"
	model2 "job-interviewer/internal/interviewer/model"
)

var (
	ErrAlreadyExistsStartedInterview = errors.New("already exists started interview")
)

type (
	CreateInterviewIn struct {
		UserID         uuid.UUID
		JobPosition    string
		JobLevel       model2.JobLevel
		QuestionsCount int64
	}

	AcceptAnswerIn struct {
		UserID uuid.UUID
		Answer string
	}

	Service interface {
		CreateInterview(ctx context.Context, in CreateInterviewIn) (*model2.Interview, error)
		StartInterview(ctx context.Context, interview *model2.Interview) error
		FinishInterview(ctx context.Context, interview *model2.Interview) error
		FindActiveInterview(ctx context.Context, userID uuid.UUID) (*model2.Interview, error)
		AcceptAnswer(ctx context.Context, in AcceptAnswerIn) (string, error)
	}
)
