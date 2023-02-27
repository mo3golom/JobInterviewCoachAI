package interview

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"job-interviewer/internal/model"
)

var (
	ErrAlreadyExistsStartedInterview = errors.New("already exists started interview")
)

type (
	StartInterviewIn struct {
		UserID         uuid.UUID
		JobPosition    string
		JobLevel       model.JobLevel
		QuestionsCount int64
	}

	AcceptAnswerIn struct {
		UserID uuid.UUID
		Answer string
	}

	Service interface {
		StartInterview(ctx context.Context, in StartInterviewIn) (*model.Interview, error)
		FinishInterview(ctx context.Context, interview *model.Interview) error
		FindActiveInterview(ctx context.Context, userID uuid.UUID) (*model.Interview, error)
		AcceptAnswer(ctx context.Context, in AcceptAnswerIn) (string, error)
	}
)
