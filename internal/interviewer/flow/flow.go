package flow

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"job-interviewer/internal/interviewer/contracts"
	"job-interviewer/internal/interviewer/flow/state"
	"job-interviewer/internal/interviewer/model"
	"job-interviewer/internal/interviewer/service/interview"
)

type DefaultInterviewFlow struct {
	startState           state.State
	waitingQuestionState state.State
	waitingAnswerState   state.State

	interviewService interview.Service
}

func NewDefaultInterviewFlow(
	interviewService interview.Service,

) *DefaultInterviewFlow {
	interviewCtx := &DefaultInterviewFlow{
		interviewService: interviewService,
	}
	interviewCtx.startState = state.NewStartState(interviewCtx)
	interviewCtx.waitingQuestionState = state.NewWaitingQuestionState(interviewCtx)
	interviewCtx.waitingAnswerState = state.NewWaitingAnswerState(interviewCtx)

	return interviewCtx
}

func (w *DefaultInterviewFlow) StartInterview(ctx context.Context, in StartInterviewIn) error {
	return w.CurrentState(nil).StartInterview(
		ctx,
		state.StartInterviewIn{
			UserID:      in.UserID,
			JobPosition: in.JobPosition,
		},
	)
}

func (w *DefaultInterviewFlow) ContinueInterview(ctx context.Context, userID uuid.UUID) error {
	activeInterview, err := w.interviewService.FindActiveInterview(ctx, userID)
	if err != nil {
		return err
	}

	return w.SetState(ctx, activeInterview.ID, model.InterviewStateWaitingQuestion)
}

func (w *DefaultInterviewFlow) FinishInterview(ctx context.Context, userID uuid.UUID) (string, error) {
	activeInterview, err := w.interviewService.FindActiveInterview(ctx, userID)
	if err != nil {
		return "", err
	}

	return w.CurrentState(activeInterview).FinishInterview(ctx, activeInterview)
}

func (w *DefaultInterviewFlow) NextQuestion(ctx context.Context, userID uuid.UUID) (*model.Question, error) {
	activeInterview, err := w.interviewService.FindActiveInterview(ctx, userID)
	if err != nil {
		return nil, err
	}

	return w.CurrentState(activeInterview).NextQuestion(ctx, activeInterview)
}

func (w *DefaultInterviewFlow) AcceptAnswer(ctx context.Context, in AcceptAnswerIn) error {
	activeInterview, err := w.interviewService.FindActiveInterview(ctx, in.UserID)
	if err != nil {
		return err
	}

	return w.CurrentState(activeInterview).AcceptAnswer(
		ctx,
		state.AcceptAnswerIn{
			Interview: activeInterview,
			Answer:    in.Answer,
		},
	)
}

func (w *DefaultInterviewFlow) GetAnswerSuggestion(ctx context.Context, userID uuid.UUID) (*model.AnswerSuggestion, error) {
	activeInterview, err := w.interviewService.FindActiveInterview(ctx, userID)
	if err != nil {
		return nil, err
	}

	return w.CurrentState(activeInterview).GetAnswerSuggestion(
		ctx,
		activeInterview,
	)
}

func (w *DefaultInterviewFlow) StartInterviewImpl(ctx context.Context, in state.StartInterviewIn) (*state.StartInterviewOut, error) {
	activeInterview, err := w.interviewService.FindActiveInterview(ctx, in.UserID)
	if err != nil && !errors.Is(err, contracts.ErrEmptyActiveInterview) {
		return nil, err
	}
	err = w.interviewService.FinishInterviewWithoutSummary(ctx, activeInterview)
	if err != nil {
		return nil, err
	}

	newInterview, err := w.interviewService.CreateInterview(
		ctx,
		interview.CreateInterviewIn{
			UserID:         in.UserID,
			JobPosition:    model.Position(in.JobPosition),
			QuestionsCount: 0,
		},
	)
	if err != nil {
		return nil, err
	}

	err = w.interviewService.StartInterview(ctx, newInterview)
	if err != nil {
		return nil, err
	}

	return &state.StartInterviewOut{
		ID: newInterview.ID,
	}, nil
}

func (w *DefaultInterviewFlow) FinishInterviewImpl(ctx context.Context, interview *model.Interview) (string, error) {
	return w.interviewService.FinishInterview(ctx, interview)
}

func (w *DefaultInterviewFlow) NextQuestionImpl(ctx context.Context, interview *model.Interview) (*model.Question, error) {
	return w.interviewService.GetNextQuestion(ctx, interview)
}

func (w *DefaultInterviewFlow) AcceptAnswerImpl(ctx context.Context, in state.AcceptAnswerIn) error {
	return w.interviewService.AcceptAnswer(
		ctx,
		interview.AcceptAnswerIn{
			Interview: in.Interview,
			Answer:    in.Answer,
		},
	)
}

func (w *DefaultInterviewFlow) GetAnswerSuggestionImpl(ctx context.Context, interview *model.Interview) (*model.AnswerSuggestion, error) {
	return w.interviewService.GetAnswerSuggestion(ctx, interview)
}

func (w *DefaultInterviewFlow) SetState(ctx context.Context, interviewID uuid.UUID, state model.InterviewState) error {
	return w.interviewService.UpdateInterviewState(ctx, interviewID, state)
}

func (w *DefaultInterviewFlow) CurrentState(interview *model.Interview) state.State {
	if interview == nil {
		return w.startState
	}

	switch interview.State {
	case model.InterviewStateStart:
		return w.startState
	case model.InterviewStateWaitingQuestion:
		return w.waitingQuestionState
	case model.InterviewStateWaitingAnswer:
		return w.waitingAnswerState
	default:
		return w.startState
	}
}
