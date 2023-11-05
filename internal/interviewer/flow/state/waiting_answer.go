package state

import (
	"context"
	"job-interviewer/internal/interviewer/model"
)

type WaitingAnswerState struct {
	interviewFlow Context
	baseState
}

func NewWaitingAnswerState(interviewFlow Context) *WaitingAnswerState {
	return &WaitingAnswerState{
		interviewFlow: interviewFlow,
	}
}

func (s *WaitingAnswerState) FinishInterview(ctx context.Context, interview *model.Interview) (string, error) {
	return s.interviewFlow.FinishInterviewImpl(ctx, interview)
}

func (s *WaitingAnswerState) AcceptAnswer(ctx context.Context, in AcceptAnswerIn) error {
	err := s.interviewFlow.AcceptAnswerImpl(ctx, in)
	if err != nil {
		return err
	}

	return s.interviewFlow.SetState(ctx, in.Interview.ID, model.InterviewStateWaitingQuestion)
}

func (s *WaitingAnswerState) GetAnswerSuggestion(ctx context.Context, interview *model.Interview) (*model.AnswerSuggestion, error) {
	return s.interviewFlow.GetAnswerSuggestionImpl(ctx, interview)
}
