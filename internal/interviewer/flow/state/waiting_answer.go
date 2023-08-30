package state

import (
	"context"
	"job-interviewer/internal/interviewer/contracts"
	"job-interviewer/internal/interviewer/model"
)

type WaitingAnswerState struct {
	interviewFlow Context
}

func NewWaitingAnswerState(interviewFlow Context) *WaitingAnswerState {
	return &WaitingAnswerState{
		interviewFlow: interviewFlow,
	}
}

func (s *WaitingAnswerState) StartInterview(_ context.Context, _ StartInterviewIn) error {
	return contracts.ErrActionDoesntAllow
}

func (s *WaitingAnswerState) FinishInterview(ctx context.Context, interview *model.Interview) (string, error) {
	return s.interviewFlow.FinishInterviewImpl(ctx, interview)
}

func (s *WaitingAnswerState) NextQuestion(_ context.Context, _ *model.Interview) (*model.Question, error) {
	return nil, contracts.ErrActionDoesntAllow
}

func (s *WaitingAnswerState) AcceptAnswer(ctx context.Context, in AcceptAnswerIn) (string, error) {
	result, err := s.interviewFlow.AcceptAnswerImpl(ctx, in)
	if err != nil {
		return "", err
	}

	err = s.interviewFlow.SetState(ctx, in.Interview.ID, model.InterviewStateWaitingQuestion)
	if err != nil {
		return "", err
	}

	return result, nil
}
