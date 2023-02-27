package question

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"job-interviewer/internal/contracts"
	"job-interviewer/internal/gpt"
	"job-interviewer/internal/model"
	"job-interviewer/internal/storage/question"
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

func (s *DefaultService) CreateQuestionsForInterview(ctx context.Context, interview *model.Interview) error {
	result, err := s.gpt.GetQuestionsList(
		ctx,
		gpt.GetQuestionsListIn{
			JobPosition:   interview.JobInfo.Position,
			QuestionCount: interview.QuestionsCount,
		},
	)
	if err != nil {
		return err
	}

	questions := convertQuestions(result, interview.JobInfo)
	return s.transactionalTemplate.Execute(ctx, func(tx transactional.Tx) error {
		err := s.storage.CreateQuestions(
			ctx,
			tx,
			questions,
		)
		if err != nil {
			return err
		}

		return s.storage.AttachQuestionsToInterview(ctx, tx, interview.ID, questions)
	})
}

func (s *DefaultService) FindNextQuestion(ctx context.Context, interviewID uuid.UUID) (*model.Question, error) {
	var result *model.Question
	err := s.transactionalTemplate.Execute(ctx, func(tx transactional.Tx) error {
		temp, err := s.storage.FindNextQuestion(ctx, tx, interviewID)
		if err != nil {
			return err
		}

		result = temp
		return nil
	})
	if errors.Is(err, question.ErrEmptyQuestionResult) {
		return nil, contracts.ErrNextQuestionEmpty
	}
	if err != nil {
		return nil, err
	}

	return result, nil
}

func convertQuestions(in []string, jobInfo model.JobInfo) []model.Question {
	result := make([]model.Question, 0, len(in))
	for _, t := range in {
		result = append(
			result,
			model.Question{
				ID:      uuid.New(),
				Text:    t,
				JobInfo: jobInfo,
			},
		)
	}

	return result
}
