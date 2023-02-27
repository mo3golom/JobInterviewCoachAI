package acceptanswer

import (
	"context"
	"job-interviewer/internal/contracts"
	"job-interviewer/internal/service/interview"
)

type UseCase struct {
	interviewService interview.Service
}

func NewUseCase(i interview.Service) *UseCase {
	return &UseCase{interviewService: i}
}

func (u *UseCase) AcceptAnswer(ctx context.Context, in contracts.AcceptAnswerIn) (string, error) {
	return u.interviewService.AcceptAnswer(
		ctx,
		interview.AcceptAnswerIn{
			UserID: in.UserID,
			Answer: in.Answer,
		},
	)
}
