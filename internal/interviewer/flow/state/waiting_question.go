package state

import (
	"context"
	"job-interviewer/internal/interviewer/contracts"
	"job-interviewer/internal/interviewer/model"
)

type WaitingQuestionState struct {
	interviewFlow Context
}

func NewWaitingQuestionState(interviewFlow Context) *WaitingQuestionState {
	return &WaitingQuestionState{
		interviewFlow: interviewFlow,
	}
}

func (s *WaitingQuestionState) StartInterview(_ context.Context, _ StartInterviewIn) error {
	return contracts.ErrActionDoesntAllow
}

func (s *WaitingQuestionState) FinishInterview(ctx context.Context, interview *model.Interview) (string, error) {
	return s.interviewFlow.FinishInterviewImpl(ctx, interview)
}

func (s *WaitingQuestionState) NextQuestion(ctx context.Context, interview *model.Interview) (*model.Question, error) {
	result, err := s.interviewFlow.NextQuestionImpl(ctx, interview)
	if err != nil {
		return nil, err
	}

	err = s.interviewFlow.SetState(ctx, interview.ID, model.InterviewStateWaitingAnswer)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *WaitingQuestionState) AcceptAnswer(_ context.Context, _ AcceptAnswerIn) (string, error) {
	return "", contracts.ErrActionDoesntAllow
}
