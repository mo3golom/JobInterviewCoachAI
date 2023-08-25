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
	defaultState             *state.Default
	waitingQuestionState     *state.WaitingQuestionState
	answeringOnQuestionState state.State

	interviewService interview.Service
}

func NewDefaultInterviewFlow(
	interviewService interview.Service,

) *DefaultInterviewFlow {
	interviewCtx := &DefaultInterviewFlow{
		interviewService: interviewService,
	}

	defaultState := state.NewDefaultState(interviewCtx)
	answeringOnQuestionState := state.NewAnsweringOnQuestionState(interviewCtx)
	waitingQuestionState := state.NewWaitingQuestionState(interviewCtx)

	interviewCtx.defaultState = defaultState
	interviewCtx.answeringOnQuestionState = answeringOnQuestionState
	interviewCtx.waitingQuestionState = waitingQuestionState

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

func (w *DefaultInterviewFlow) AcceptAnswer(ctx context.Context, in AcceptAnswerIn) (string, error) {
	activeInterview, err := w.interviewService.FindActiveInterview(ctx, in.UserID)
	if err != nil {
		return "", err
	}

	return w.CurrentState(activeInterview).AcceptAnswer(
		ctx,
		state.AcceptAnswerIn{
			Interview: activeInterview,
			Answer:    in.Answer,
		},
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
			JobPosition:    in.JobPosition,
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

func (w *DefaultInterviewFlow) AcceptAnswerImpl(ctx context.Context, in state.AcceptAnswerIn) (string, error) {
	return w.interviewService.AcceptAnswer(
		ctx,
		interview.AcceptAnswerIn{
			Interview: in.Interview,
			Answer:    in.Answer,
		},
	)
}

func (w *DefaultInterviewFlow) SetState(ctx context.Context, interviewID uuid.UUID, state model.InterviewState) error {
	return w.interviewService.UpdateInterviewState(ctx, interviewID, state)
}

func (w *DefaultInterviewFlow) CurrentState(interview *model.Interview) state.State {
	if interview == nil {
		return w.defaultState
	}

	switch interview.State {
	case model.InterviewStateDefault:
		return w.defaultState
	case model.InterviewStateWaitingQuestion:
		return w.waitingQuestionState
	case model.InterviewStateAnsweringOnQuestion:
		return w.answeringOnQuestionState
	default:
		return w.defaultState
	}
}
