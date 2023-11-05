package state

import (
	"context"
	"github.com/google/uuid"
	"job-interviewer/internal/interviewer/contracts"
	"job-interviewer/internal/interviewer/model"
)

type baseState struct {
}

func (b *baseState) StartInterview(_ context.Context, _ StartInterviewIn) error {
	return contracts.ErrActionDoesntAllow
}

func (b *baseState) FinishInterview(_ context.Context, _ *model.Interview) (string, error) {
	return "", contracts.ErrActionDoesntAllow
}

func (b *baseState) NextQuestion(_ context.Context, _ *model.Interview) (*model.Question, error) {
	return nil, contracts.ErrActionDoesntAllow
}

func (b *baseState) AcceptAnswer(_ context.Context, _ AcceptAnswerIn) error {
	return contracts.ErrActionDoesntAllow
}

func (b *baseState) ContinueInterview(_ context.Context, _ uuid.UUID) error {
	return contracts.ErrActionDoesntAllow
}

func (b *baseState) GetAnswerSuggestion(_ context.Context, _ *model.Interview) (*model.AnswerSuggestion, error) {
	return nil, contracts.ErrActionDoesntAllow
}
