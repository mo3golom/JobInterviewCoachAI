package interview

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"job-interviewer/internal/interviewer/model"
	"job-interviewer/pkg/transactional"
)

var (
	ErrEmptyInterviewResult = errors.New("empty interview result")
)

type Storage interface {
	CreateInterview(ctx context.Context, tx transactional.Tx, interview *model.Interview) error
	UpdateInterview(ctx context.Context, tx transactional.Tx, interview *model.Interview) error
	UpdateInterviewState(ctx context.Context, tx transactional.Tx, interviewID uuid.UUID, state model.InterviewState) error
	FindActiveInterviewByUserID(ctx context.Context, tx transactional.Tx, userID uuid.UUID) (*model.Interview, error)
}
