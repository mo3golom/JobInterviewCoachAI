package question

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"job-interviewer/internal/interviewer/gpt"
	"job-interviewer/internal/interviewer/model"
	"job-interviewer/internal/interviewer/storage/question"
	"job-interviewer/pkg/transactional"
)

type DefaultService struct {
	gpt                   gpt.Gateway
	storage               question.Storage
	transactionalTemplate transactional.Template
}

func NewQuestionService(g gpt.Gateway, s question.Storage, tr transactional.Template) *DefaultService {
	return &DefaultService{
		gpt:                   g,
		storage:               s,
		transactionalTemplate: tr,
	}
}

func (s *DefaultService) GetNextQuestion(ctx context.Context, interview *model.Interview) (*model.Question, error) {
	var nextQuestion *model.Question
	err := s.transactionalTemplate.Execute(ctx, func(tx transactional.Tx) error {
		tempQuestion, err := s.storage.FindActiveQuestionByInterviewID(ctx, tx, interview.ID)
		if err != nil {
			return err
		}

		nextQuestion = tempQuestion
		return nil
	})
	if err != nil && !errors.Is(err, question.ErrEmptyQuestionResult) {
		return nil, err
	}

	if nextQuestion != nil {
		return nextQuestion, nil
	}

	result, err := s.gpt.GetQuestion(ctx, interview.JobInfo.Position)
	if err != nil {
		return nil, err
	}

	nextQuestion = &model.Question{
		ID:      uuid.New(),
		Text:    result,
		JobInfo: interview.JobInfo,
	}
	err = s.transactionalTemplate.Execute(ctx, func(tx transactional.Tx) error {
		err := s.storage.CreateQuestion(
			ctx,
			tx,
			nextQuestion,
		)
		if err != nil {
			return err
		}

		return s.storage.AttachQuestionToInterview(ctx, tx, interview.ID, nextQuestion)
	})
	if err != nil {
		return nil, err
	}

	return nextQuestion, nil
}
