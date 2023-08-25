package acceptanswer

import (
	"context"
	"job-interviewer/internal/interviewer/contracts"
	"job-interviewer/internal/interviewer/flow"
)

type UseCase struct {
	interviewFlow flow.InterviewFlow
}

func NewUseCase(interviewFlow flow.InterviewFlow) *UseCase {
	return &UseCase{
		interviewFlow: interviewFlow,
	}
}

func (u *UseCase) AcceptAnswer(ctx context.Context, in contracts.AcceptAnswerIn) (string, error) {
	return u.interviewFlow.AcceptAnswer(
		ctx,
		flow.AcceptAnswerIn{
			UserID: in.UserID,
			Answer: in.Answer,
		},
	)
}
