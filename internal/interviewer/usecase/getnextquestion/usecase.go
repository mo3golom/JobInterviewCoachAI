package getnextquestion

import (
	"context"
	"github.com/google/uuid"
	"job-interviewer/internal/interviewer/model"
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

func (u *UseCase) GetNextQuestion(ctx context.Context, userID uuid.UUID) (*model.Question, error) {
	activeInterview, err := u.interviewService.FindActiveInterview(ctx, userID)
	if err != nil {
		return nil, err
	}

	nextQuestion, err := u.questionService.GetNextQuestion(ctx, activeInterview)
	if err != nil {
		return nil, err
	}

	return nextQuestion, nil
}
