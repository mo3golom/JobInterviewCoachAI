package state

import (
	"context"
	"job-interviewer/internal/interviewer/model"
)

type StartFlow struct {
	interviewFlow Context
	baseState
}

func NewStartState(interviewFlow Context) *StartFlow {
	return &StartFlow{
		interviewFlow: interviewFlow,
	}
}

func (s *StartFlow) StartInterview(ctx context.Context, in StartInterviewIn) error {
	result, err := s.interviewFlow.StartInterviewImpl(ctx, in)
	if err != nil {
		return nil
	}

	return s.interviewFlow.SetState(ctx, result.ID, model.InterviewStateWaitingQuestion)
}
