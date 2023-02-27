package getnextquestion

import (
	"context"
	"github.com/google/uuid"
	"job-interviewer/internal/model"
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

func (u *UseCase) GetNextQuestion(ctx context.Context, userID uuid.UUID) (*model.Question, error) {
	activeInterview, err := u.interviewService.FindActiveInterview(ctx, userID)
	if err != nil {
		return nil, err
	}

	nextQuestion, err := u.questionService.FindNextQuestion(ctx, activeInterview.ID)
	if err != nil {
		return nil, err
	}

	return nextQuestion, nil
}
