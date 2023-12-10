package flow

import (
	"context"
	"github.com/google/uuid"
	"job-interviewer/internal"
	"job-interviewer/internal/interviewer/model"
)

type (
	StartInterviewIn struct {
		UserID      uuid.UUID
		JobPosition internal.Position
	}

	StartInterviewOut struct {
		ID uuid.UUID
	}

	AcceptAnswerIn struct {
		UserID uuid.UUID
		Answer string
	}

	InterviewFlow interface {
		StartInterview(ctx context.Context, in StartInterviewIn) error
		FinishInterview(ctx context.Context, userID uuid.UUID) (string, error)
		NextQuestion(ctx context.Context, userID uuid.UUID) (*model.Question, error)
		AcceptAnswer(ctx context.Context, in AcceptAnswerIn) error
		ContinueInterview(ctx context.Context, userID uuid.UUID) error
		GetAnswerSuggestion(ctx context.Context, userID uuid.UUID) (*model.AnswerSuggestion, error)
	}
)
