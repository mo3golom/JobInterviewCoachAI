package startinterview

import (
	"context"
	"job-interviewer/internal/interviewer/contracts"
	"job-interviewer/internal/interviewer/flow"
	"job-interviewer/internal/interviewer/service/interview"
)

type UseCase struct {
	interviewService interview.Service
	interviewFlow    flow.InterviewFlow
}

func NewUseCase(i interview.Service, interviewFlow flow.InterviewFlow) *UseCase {
	return &UseCase{
		interviewService: i,
		interviewFlow:    interviewFlow,
	}
}

func (u *UseCase) StartInterview(ctx context.Context, in contracts.StartInterviewIn) error {
	return u.interviewFlow.StartInterview(
		ctx,
		flow.StartInterviewIn{
			UserID:      in.UserID,
			JobPosition: in.Questions.JobPosition,
		},
	)
}
