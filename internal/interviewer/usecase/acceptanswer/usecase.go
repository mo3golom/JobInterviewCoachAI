package acceptanswer

import (
	"context"
	"github.com/google/uuid"
	"job-interviewer/internal/interviewer/contracts"
	"job-interviewer/internal/interviewer/flow"
	"job-interviewer/internal/interviewer/model"
)

const (
	skipAnswerText = "I don't want to answer this question, skip the question and continue"
)

type UseCase struct {
	interviewFlow flow.InterviewFlow
}

func NewUseCase(interviewFlow flow.InterviewFlow) *UseCase {
	return &UseCase{
		interviewFlow: interviewFlow,
	}
}

func (u *UseCase) AcceptAnswer(ctx context.Context, in contracts.AcceptAnswerIn) error {
	return u.interviewFlow.AcceptAnswer(
		ctx,
		flow.AcceptAnswerIn{
			UserID: in.UserID,
			Answer: in.Answer,
		},
	)
}

func (u *UseCase) SkipQuestion(ctx context.Context, userID uuid.UUID) error {
	return u.interviewFlow.AcceptAnswer(
		ctx,
		flow.AcceptAnswerIn{
			UserID: userID,
			Answer: skipAnswerText,
		},
	)
}

func (u *UseCase) GetAnswerSuggestion(ctx context.Context, userID uuid.UUID) (*model.AnswerSuggestion, error) {
	return u.interviewFlow.GetAnswerSuggestion(ctx, userID)
}
