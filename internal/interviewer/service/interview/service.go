package interview

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"job-interviewer/internal"
	"job-interviewer/internal/interviewer/contracts"
	"job-interviewer/internal/interviewer/gpt"
	"job-interviewer/internal/interviewer/model"
	"job-interviewer/internal/interviewer/service/subscription"
	"job-interviewer/internal/interviewer/storage/interview"
	"job-interviewer/internal/interviewer/storage/messages"
	"job-interviewer/pkg/transactional"
)

const (
	startTechnicalInterviewPrompt  = `I want you to act as an interviewer. I will be the candidate and you will ask me the technical interview questions for the %s position. I want you to only reply as the interviewer. Do not write all the conservation at once. I want you to only do the interview with me. Ask me the tricky questions and wait for my answers. Do not write explanations. Do not write "interviewer:". Ask me the questions one by one like an interviewer does and wait for my answers. My first sentence is "Hi"`
	startBehavioralInterviewPrompt = `I want you to act as an interviewer. I will be the candidate and you will ask me the behavioral interview questions. I want you to only reply as the interviewer. Do not write all the conservation at once. I want you to only do the interview with me. Ask me the tricky questions and wait for my answers. Do not write explanations. Do not write "interviewer:". Ask me the questions one by one like an interviewer does and wait for my answers. My first sentence is "Hi"`
)

type DefaultService struct {
	gpt                   gpt.Gateway
	interviewStorage      interview.Storage
	messagesStorage       messages.Storage
	subscriptionService   subscription.Service
	transactionalTemplate transactional.Template
}

func NewInterviewService(
	g gpt.Gateway,
	is interview.Storage,
	messagesStorage messages.Storage,
	tr transactional.Template,
	subscriptionService subscription.Service,
) *DefaultService {
	return &DefaultService{
		gpt:                   g,
		interviewStorage:      is,
		messagesStorage:       messagesStorage,
		transactionalTemplate: tr,
		subscriptionService:   subscriptionService,
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

		return s.subscriptionService.DecreaseFreeAttempts(ctx, tx, interview.UserID)
	})
}

func (s *DefaultService) FinishInterview(ctx context.Context, interview *model.Interview) (string, error) {
	history, err := s.messagesStorage.GetMessagesByInterviewID(ctx, interview.ID)
	if err != nil {
		return "", err
	}

	var summary string
	if len(history) > 2 {
		result, err := s.gpt.SummarizeDialogue(ctx, history)
		if err != nil {
			return "", err
		}

		summary = result.Content
	}

	err = s.FinishInterviewWithoutSummary(ctx, interview)
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

func (s *DefaultService) GetNextQuestion(ctx context.Context, interview *model.Interview) (*model.Question, error) {
	history, err := s.messagesStorage.GetMessagesByInterviewID(ctx, interview.ID)
	if err != nil {
		return nil, err
	}

	if len(history) == 0 {
		startPrompt := fmt.Sprintf(startTechnicalInterviewPrompt, interview.JobInfo.Position)
		if interview.JobInfo.Position == internal.Behavioral {
			startPrompt = startBehavioralInterviewPrompt
		}

		result, err := s.gpt.StartDialogue(
			ctx,
			startPrompt,
		)
		if err != nil {
			return nil, err
		}

		err = s.transactionalTemplate.Execute(ctx, func(tx transactional.Tx) error {
			return s.messagesStorage.CreateMessages(
				ctx,
				tx,
				interview.ID,
				[]model.Message{
					{
						Role:    model.RoleUser,
						Content: startPrompt,
					},
					{
						Role:    model.RoleAssistant,
						Content: result.Content,
					},
				},
			)
		})
		if err != nil {
			return nil, err
		}

		return &model.Question{
			Text: result.Content,
		}, nil
	}

	if lastMessage := history[len(history)-1]; lastMessage.Role == model.RoleAssistant {
		return &model.Question{
			Text: lastMessage.Content,
		}, nil
	}

	result, err := s.gpt.ContinueDialogue(ctx, history)
	if err != nil {
		return nil, err
	}

	return &model.Question{
		Text: result.Content,
	}, nil
}

func (s *DefaultService) AcceptAnswer(ctx context.Context, in AcceptAnswerIn) error {
	history, err := s.messagesStorage.GetMessagesByInterviewID(ctx, in.Interview.ID)
	if err != nil {
		return err
	}

	history = append(
		history,
		model.Message{
			Role:    model.RoleUser,
			Content: in.Answer,
		},
	)

	result, err := s.gpt.ContinueDialogue(ctx, history)
	if err != nil {
		return err
	}

	return s.transactionalTemplate.Execute(ctx, func(tx transactional.Tx) error {
		err := s.messagesStorage.CreateMessage(
			ctx,
			tx,
			in.Interview.ID,
			&model.Message{
				Role:    model.RoleUser,
				Content: in.Answer,
			},
		)
		if err != nil {
			return err
		}

		return s.messagesStorage.CreateMessage(
			ctx,
			tx,
			in.Interview.ID,
			&model.Message{
				Role:    model.RoleAssistant,
				Content: result.Content,
			},
		)
	})
}

func (s *DefaultService) GetAnswerSuggestion(ctx context.Context, interview *model.Interview) (*model.AnswerSuggestion, error) {
	history, err := s.messagesStorage.GetMessagesByInterviewID(ctx, interview.ID)
	if err != nil {
		return nil, err
	}

	if len(history) == 0 {
		return nil, contracts.ErrInterviewQuestionsIsEmpty
	}

	if lastMessage := history[len(history)-1]; lastMessage.Role == model.RoleUser {
		return nil, contracts.ErrInterviewQuestionsIsEmpty
	}

	result, err := s.gpt.GetAnswerSuggestion(ctx, history)
	if err != nil {
		return nil, err
	}

	return &model.AnswerSuggestion{
		Text: result.Content,
	}, nil
}

func (s *DefaultService) UpdateInterviewState(ctx context.Context, interviewID uuid.UUID, state model.InterviewState) error {
	return s.transactionalTemplate.Execute(ctx, func(tx transactional.Tx) error {
		return s.interviewStorage.UpdateInterviewState(
			ctx,
			tx,
			interviewID,
			state,
		)
	})
}
