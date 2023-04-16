package interview

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"job-interviewer/internal/interviewer/contracts"
	"job-interviewer/internal/interviewer/gpt"
	"job-interviewer/internal/interviewer/model"
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

func (s *DefaultService) CreateInterview(ctx context.Context, in CreateInterviewIn) (*model.Interview, error) {
	newInterview := &model.Interview{
		ID:     uuid.New(),
		UserID: in.UserID,
		Status: model.InterviewStatusCreated,
		JobInfo: model.JobInfo{
			Position: in.JobPosition,
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

func (s *DefaultService) StartInterview(ctx context.Context, interview *model.Interview) error {
	interview.Status = model.InterviewStatusStarted
	return s.transactionalTemplate.Execute(ctx, func(tx transactional.Tx) error {
		return s.interviewStorage.UpdateInterview(
			ctx,
			tx,
			interview,
		)
	})
}

func (s *DefaultService) FinishInterviewWithoutSummary(ctx context.Context, interview *model.Interview) error {
	if interview == nil {
		return nil
	}

	interview.Status = model.InterviewStatusFinished
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
				Current:     model.InterviewQuestionStatusActive,
				Target:      model.InterviewQuestionStatusCanceled,
			},
		)
	})
}

func (s *DefaultService) FinishInterview(ctx context.Context, interview *model.Interview) (string, error) {
	err := s.FinishInterviewWithoutSummary(ctx, interview)
	if err != nil {
		return "", err
	}

	answersComments, err := s.questionStorage.FindAnswersCommentsByInterviewID(ctx, interview.ID)
	if err != nil {
		return "", err
	}
	if len(answersComments) == 0 {
		return "", nil
	}

	summary, err := s.gpt.SummarizeAnswersComments(ctx, answersComments)
	if err != nil {
		return "", err
	}

	return summary, nil
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
	var activeInterviewID uuid.UUID
	err := s.transactionalTemplate.Execute(ctx, func(tx transactional.Tx) error {
		activeInterview, err := s.interviewStorage.FindActiveInterviewByUserID(ctx, tx, in.UserID)
		if errors.Is(err, interview.ErrEmptyInterviewResult) {
			return contracts.ErrEmptyActiveInterview
		}
		if err != nil {
			return err
		}

		tempQuestion, err := s.questionStorage.FindActiveQuestionByInterviewID(ctx, tx, activeInterview.ID)
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
			Answer:   in.Answer,
			Question: currentQuestion.Text,
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
