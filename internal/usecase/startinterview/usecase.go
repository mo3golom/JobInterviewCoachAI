package startinterview

import (
	"context"
	"job-interviewer/internal/contracts"
	"job-interviewer/internal/service/interview"
	"job-interviewer/internal/service/question"
)

type UseCase struct {
	interviewService interview.Service
	questionService  question.Service
}

func NewUseCase(i interview.Service, q question.Service) *UseCase {
	return &UseCase{interviewService: i, questionService: q}
}

func (u *UseCase) StartInterview(ctx context.Context, in contracts.StartInterviewIn) error {
	newInterview, err := u.interviewService.StartInterview(
		ctx,
		interview.StartInterviewIn{
			UserID:         in.UserID,
			JobPosition:    in.JobPosition,
			JobLevel:       in.JobLevel,
			QuestionsCount: in.QuestionsCount,
		},
	)
	if err != nil {
		return err
	}

	return u.questionService.CreateQuestionsForInterview(ctx, newInterview)
}
