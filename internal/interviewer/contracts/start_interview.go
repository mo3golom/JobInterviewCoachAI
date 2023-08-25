package contracts

import (
	"context"
	"github.com/google/uuid"
)

type (
	StartInterviewIn struct {
		UserID    uuid.UUID
		Questions StartInterviewQuestionsIn
	}

	StartInterviewQuestionsIn struct {
		JobPosition string
	}

	StartInterviewUseCase interface {
		StartInterview(ctx context.Context, in StartInterviewIn) error
	}
)
