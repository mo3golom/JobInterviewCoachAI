package startinterview

import (
	"context"
	"errors"
	"job-interviewer/internal/interviewer/contracts"
	"job-interviewer/internal/interviewer/service/interview"
	"job-interviewer/internal/interviewer/service/question"
)

type UseCase struct {
	interviewService interview.Service
	questionService  question.Service
}

func NewUseCase(i interview.Service, q question.Service) *UseCase {
	return &UseCase{interviewService: i, questionService: q}
}

func (u *UseCase) StartInterview(ctx context.Context, in contracts.StartInterviewIn) error {
	activeInterview, err := u.interviewService.FindActiveInterview(ctx, in.UserID)
	if err != nil && !errors.Is(err, contracts.ErrEmptyActiveInterview) {
		return err
	}
	err = u.interviewService.FinishInterviewWithoutSummary(ctx, activeInterview)
	if err != nil {
		return err
	}

	newInterview, err := u.interviewService.CreateInterview(
		ctx,
		interview.CreateInterviewIn{
			UserID:         in.UserID,
			JobPosition:    in.JobPosition,
			QuestionsCount: in.QuestionsCount,
		},
	)
	if err != nil {
		return err
	}

	return u.interviewService.StartInterview(ctx, newInterview)
}
