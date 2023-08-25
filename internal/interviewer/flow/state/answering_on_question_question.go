package state

import (
	"context"
	"job-interviewer/internal/interviewer/model"
)

type AnsweringOnQuestionState struct {
	interviewFlow Context
}

func NewAnsweringOnQuestionState(interviewFlow Context) *AnsweringOnQuestionState {
	return &AnsweringOnQuestionState{
		interviewFlow: interviewFlow,
	}
}

func (s *AnsweringOnQuestionState) StartInterview(_ context.Context, _ StartInterviewIn) error {
	return nil
}

func (s *AnsweringOnQuestionState) FinishInterview(ctx context.Context, interview *model.Interview) (string, error) {
	return s.interviewFlow.FinishInterviewImpl(ctx, interview)
}

func (s *AnsweringOnQuestionState) NextQuestion(ctx context.Context, interview *model.Interview) (*model.Question, error) {
	result, err := s.interviewFlow.NextQuestionImpl(ctx, interview)
	if err != nil {
		return nil, err
	}

	err = s.interviewFlow.SetState(ctx, interview.ID, model.InterviewStateWaitingQuestion)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *AnsweringOnQuestionState) AcceptAnswer(ctx context.Context, in AcceptAnswerIn) (string, error) {
	return s.interviewFlow.AcceptAnswerImpl(ctx, in)
}
