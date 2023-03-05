package interview

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"job-interviewer/internal/interviewer/contracts"
	"job-interviewer/internal/interviewer/gpt"
	model2 "job-interviewer/internal/interviewer/model"
	"job-interviewer/internal/interviewer/storage/interview"
	"job-interviewer/internal/interviewer/storage/question"
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

func (s *DefaultService) CreateInterview(ctx context.Context, in CreateInterviewIn) (*model2.Interview, error) {
	newInterview := &model2.Interview{
		ID:     uuid.New(),
		UserID: in.UserID,
		Status: model2.InterviewStatusCreated,
		JobInfo: model2.JobInfo{
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

func (s *DefaultService) StartInterview(ctx context.Context, interview *model2.Interview) error {
	interview.Status = model2.InterviewStatusStarted
	return s.transactionalTemplate.Execute(ctx, func(tx transactional.Tx) error {
		return s.interviewStorage.UpdateInterview(
			ctx,
			tx,
			interview,
		)
	})
}

func (s *DefaultService) FinishInterview(ctx context.Context, interview *model2.Interview) error {
	if interview == nil {
		return nil
	}

	interview.Status = model2.InterviewStatusFinished
	return s.transactionalTemplate.Execute(ctx, func(tx transactional.Tx) error {
		err := s.interviewStorage.UpdateInterview(ctx, tx, interview)
		if err != nil {
			return err
		}

		return s.questionStorage.UpdateInterviewQuestionStatus(
			ctx,
			tx,
			question.UpdateInterviewQuestionStatusIn{
				InterviewID: interview.ID,
				Current:     model2.InterviewQuestionStatusCreated,
				Target:      model2.InterviewQuestionStatusCanceled,
			},
		)
	})
}

func (s *DefaultService) FindActiveInterview(ctx context.Context, userID uuid.UUID) (*model2.Interview, error) {
	var existingInterview *model2.Interview
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
	var currentQuestion *model2.Question
	var activeInterviewID uuid.UUID
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
		activeInterviewID = activeInterview.ID

		return nil
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

	err = s.transactionalTemplate.Execute(ctx, func(tx transactional.Tx) error {
		return s.questionStorage.SetQuestionAnswered(
			ctx,
			tx,
			question.SetQuestionAnsweredIn{
				InterviewID: activeInterviewID,
				QuestionID:  currentQuestion.ID,
				Answer:      in.Answer,
				GptComment:  out,
			},
		)
	})
	if err != nil {
		return "", err
	}

	return out, nil
}
