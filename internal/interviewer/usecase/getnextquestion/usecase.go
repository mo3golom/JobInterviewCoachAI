package getnextquestion

import (
	"context"
	"github.com/google/uuid"
	"job-interviewer/internal/interviewer/flow"
	"job-interviewer/internal/interviewer/model"
)

type UseCase struct {
	interviewFlow flow.InterviewFlow
}

func NewUseCase(interviewFlow flow.InterviewFlow) *UseCase {
	return &UseCase{
		interviewFlow: interviewFlow,
	}
}

func (u *UseCase) GetNextQuestion(ctx context.Context, userID uuid.UUID) (*model.Question, error) {
	return u.interviewFlow.NextQuestion(ctx, userID)
}
