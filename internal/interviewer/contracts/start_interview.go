package contracts

import (
	"context"
	"github.com/google/uuid"
	"job-interviewer/internal/interviewer/model"
)

type (
	StartInterviewIn struct {
		UserID    uuid.UUID
		Questions StartInterviewQuestionsIn
	}

	StartInterviewQuestionsIn struct {
		JobPosition model.Position
	}

	StartInterviewUseCase interface {
		StartInterview(ctx context.Context, in StartInterviewIn) error
		ContinueInterview(ctx context.Context, userID uuid.UUID) error
	}
)
