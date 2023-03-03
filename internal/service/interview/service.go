package interview

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"job-interviewer/internal/contracts"
	"job-interviewer/internal/gpt"
	"job-interviewer/internal/model"
	"job-interviewer/internal/storage/interview"
	"job-interviewer/internal/storage/question"
	"job-interviewer/pkg/transactional"
)

type DefaultService struct {
	gpt                   gpt.Gateway
	interviewStorage      interview.Storage
	questionStorage       question.Storage
	transactionalTemplate transactional.Template
}

func NewInterviewService(
	g gpt.Gateway,
	is interview.Storage,
	qs question.Storage,
	tr transactional.Template,
) *DefaultService {
	return &DefaultService{
		gpt:                   g,
		interviewStorage:      is,
		questionStorage:       qs,
		transactionalTemplate: tr,
	}
}

func (s *DefaultService) StartInterview(ctx context.Context, in StartInterviewIn) (*model.Interview, error) {
	newInterview := &model.Interview{
		ID:     uuid.New(),
		UserID: in.UserID,
		JobInfo: model.JobInfo{
			Position: in.JobPosition,
			Level:    in.JobLevel,
		},
		QuestionsCount: in.QuestionsCount,
	}

	err := s.transactionalTemplate.Execute(ctx, func(tx transactional.Tx) error {
		existing, err := s.interviewStorage.FindActiveInterviewByUserID(ctx, tx, in.UserID)
		if err != nil && !errors.Is(err, interview.ErrEmptyInterviewResult) {
			return err
		}
		if existing != nil {
			return ErrAlreadyExistsStartedInterview
		}

		return s.interviewStorage.CreateInterview(
			ctx,
			tx,
			newInterview,
		)
	})
	if err != nil {
		return nil, err
	}

	return newInterview, nil
}

func (s *DefaultService) FinishInterview(ctx context.Context, interview *model.Interview) error {
	if interview == nil {
		return nil
	}

	interview.Status = model.InterviewStatusFinished
	return s.transactionalTemplate.Execute(ctx, func(tx transactional.Tx) error {
		return s.interviewStorage.UpdateInterview(ctx, tx, interview)
	})
}

func (s *DefaultService) FindActiveInterview(ctx context.Context, userID uuid.UUID) (*model.Interview, error) {
	var existingInterview *model.Interview
	err := s.transactionalTemplate.Execute(ctx, func(tx transactional.Tx) error {
		temp, err := s.interviewStorage.FindActiveInterviewByUserID(ctx, tx, userID)
		if err != nil {
			return err
		}

		existingInterview = temp
		return nil
	})
	if errors.Is(err, interview.ErrEmptyInterviewResult) {
		return nil, contracts.ErrEmptyActiveInterview
	}
	if err != nil {
		return nil, err
	}

	return existingInterview, nil
}

func (s *DefaultService) AcceptAnswer(ctx context.Context, in AcceptAnswerIn) (string, error) {
	var currentQuestion *model.Question
	err := s.transactionalTemplate.Execute(ctx, func(tx transactional.Tx) error {
		activeInterview, err := s.interviewStorage.FindActiveInterviewByUserID(ctx, tx, in.UserID)
		if errors.Is(err, interview.ErrEmptyInterviewResult) {
			return contracts.ErrEmptyActiveInterview
		}
		if err != nil {
			return err
		}

		tempQuestion, err := s.questionStorage.FindActiveQuestionByInterviewID(ctx, activeInterview.ID)
		if errors.Is(err, question.ErrEmptyQuestionResult) {
			return contracts.ErrInterviewQuestionsIsEmpty
		}
		if err != nil {
			return err
		}
		currentQuestion = tempQuestion

		return s.questionStorage.SetQuestionAnswered(
			ctx,
			tx,
			activeInterview.ID,
			tempQuestion.ID,
		)
	})
	if err != nil {
		return "", err
	}

	out, err := s.gpt.AcceptAnswer(
		ctx,
		gpt.AcceptAnswerIn{
			Answer:      in.Answer,
			Question:    currentQuestion.Text,
			JobPosition: currentQuestion.JobInfo.Position,
		},
	)
	if err != nil {
		return "", err
	}

	return out, nil
}
