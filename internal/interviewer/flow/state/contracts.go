package state

import (
	"context"
	"github.com/google/uuid"
	"job-interviewer/internal/interviewer/model"
)

type (
	StartInterviewIn struct {
		UserID      uuid.UUID
		JobPosition string
	}

	StartInterviewOut struct {
		ID uuid.UUID
	}

	AcceptAnswerIn struct {
		Interview *model.Interview
		Answer    string
	}

	Context interface {
		StartInterviewImpl(ctx context.Context, in StartInterviewIn) (*StartInterviewOut, error)
		FinishInterviewImpl(ctx context.Context, interview *model.Interview) (string, error)
		NextQuestionImpl(ctx context.Context, interview *model.Interview) (*model.Question, error)
		AcceptAnswerImpl(ctx context.Context, in AcceptAnswerIn) error
		GetAnswerSuggestionImpl(ctx context.Context, interview *model.Interview) (*model.AnswerSuggestion, error)
		SetState(ctx context.Context, interviewID uuid.UUID, state model.InterviewState) error
		CurrentState(interview *model.Interview) State
	}

	State interface {
		StartInterview(ctx context.Context, in StartInterviewIn) error
		FinishInterview(ctx context.Context, interview *model.Interview) (string, error)
		NextQuestion(ctx context.Context, interview *model.Interview) (*model.Question, error)
		AcceptAnswer(ctx context.Context, in AcceptAnswerIn) error
		GetAnswerSuggestion(ctx context.Context, interview *model.Interview) (*model.AnswerSuggestion, error)
	}
)
