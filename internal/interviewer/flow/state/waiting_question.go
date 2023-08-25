package state

import (
	"context"
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

func (w *WaitingQuestionState) StartInterview(_ context.Context, _ StartInterviewIn) error {
	return nil
}

func (w *WaitingQuestionState) FinishInterview(ctx context.Context, interview *model.Interview) (string, error) {
	return w.interviewFlow.FinishInterviewImpl(ctx, interview)
}

func (w *WaitingQuestionState) NextQuestion(ctx context.Context, interview *model.Interview) (*model.Question, error) {
	return w.interviewFlow.NextQuestionImpl(ctx, interview)
}

func (w *WaitingQuestionState) AcceptAnswer(ctx context.Context, in AcceptAnswerIn) (string, error) {
	result, err := w.interviewFlow.AcceptAnswerImpl(ctx, in)
	if err != nil {
		return "", err
	}

	err = w.interviewFlow.SetState(ctx, in.Interview.ID, model.InterviewStateAnsweringOnQuestion)
	if err != nil {
		return "", err
	}

	return result, nil
}
