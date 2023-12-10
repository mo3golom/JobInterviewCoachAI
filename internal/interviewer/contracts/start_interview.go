package contracts

import (
	"context"
	"github.com/google/uuid"
	"job-interviewer/internal"
)

type (
	StartInterviewIn struct {
		UserID    uuid.UUID
		Questions StartInterviewQuestionsIn
	}

	StartInterviewQuestionsIn struct {
		JobPosition internal.Position
	}

	StartInterviewUseCase interface {
		StartInterview(ctx context.Context, in StartInterviewIn) error
		ContinueInterview(ctx context.Context, userID uuid.UUID) error
	}
)
