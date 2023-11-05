package messages

import (
	"context"
	"github.com/google/uuid"
	"job-interviewer/internal/interviewer/model"
	"job-interviewer/pkg/transactional"
)

type (
	Storage interface {
		CreateMessage(ctx context.Context, tx transactional.Tx, interviewID uuid.UUID, in *model.Message) error
		GetMessagesByInterviewID(ctx context.Context, interviewID uuid.UUID) ([]model.Message, error)
		GetMessagesFromActiveInterviewByUserID(ctx context.Context, userID uuid.UUID) ([]model.Message, error)
	}
)
