package state

import (
	"context"
	"job-interviewer/internal/interviewer/contracts"
	"job-interviewer/internal/interviewer/model"
)

type Default struct {
	interviewFlow Context
}

func NewDefaultState(interviewFlow Context) *Default {
	return &Default{
		interviewFlow: interviewFlow,
	}
}

func (s *Default) StartInterview(ctx context.Context, in StartInterviewIn) error {
	result, err := s.interviewFlow.StartInterviewImpl(ctx, in)
	if err != nil {
		return nil
	}

	return s.interviewFlow.SetState(ctx, result.ID, model.InterviewStateWaitingQuestion)
}

func (s *Default) FinishInterview(_ context.Context, _ *model.Interview) (string, error) {
	return "", contracts.ErrEmptyActiveInterview
}

func (s *Default) NextQuestion(_ context.Context, _ *model.Interview) (*model.Question, error) {
	return nil, contracts.ErrEmptyActiveInterview
}

func (s *Default) AcceptAnswer(_ context.Context, _ AcceptAnswerIn) (string, error) {
	return "", contracts.ErrEmptyActiveInterview
}
