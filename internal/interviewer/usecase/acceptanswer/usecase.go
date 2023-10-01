package acceptanswer

import (
	"context"
	"github.com/google/uuid"
	"job-interviewer/internal/interviewer/contracts"
	"job-interviewer/internal/interviewer/flow"
	"job-interviewer/internal/interviewer/model"
	"job-interviewer/internal/interviewer/service/subscription"
)

const (
	skipAnswerText = "I don't want to answer this question, skip the question and continue"
)

type UseCase struct {
	interviewFlow       flow.InterviewFlow
	subscriptionService subscription.Service
}

func NewUseCase(
	interviewFlow flow.InterviewFlow,
	subscriptionService subscription.Service,
) *UseCase {
	return &UseCase{
		interviewFlow:       interviewFlow,
		subscriptionService: subscriptionService,
	}
}

func (u *UseCase) AcceptAnswer(ctx context.Context, in contracts.AcceptAnswerIn) error {
	available, err := u.subscriptionService.IsAvailable(ctx, in.UserID)
	if err != nil {
		return err
	}
	if !available.Result {
		return available.Reason
	}

	return u.interviewFlow.AcceptAnswer(
		ctx,
		flow.AcceptAnswerIn{
			UserID: in.UserID,
			Answer: in.Answer,
		},
	)
}

func (u *UseCase) SkipQuestion(ctx context.Context, userID uuid.UUID) error {
	available, err := u.subscriptionService.IsAvailable(ctx, userID)
	if err != nil {
		return err
	}
	if !available.Result {
		return available.Reason
	}

	return u.interviewFlow.AcceptAnswer(
		ctx,
		flow.AcceptAnswerIn{
			UserID: userID,
			Answer: skipAnswerText,
		},
	)
}

func (u *UseCase) GetAnswerSuggestion(ctx context.Context, userID uuid.UUID) (*model.AnswerSuggestion, error) {
	available, err := u.subscriptionService.IsAvailable(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !available.Result {
		return nil, available.Reason
	}

	return u.interviewFlow.GetAnswerSuggestion(ctx, userID)
}
