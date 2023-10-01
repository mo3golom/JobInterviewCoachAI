package subscription

import (
	"context"
	"github.com/google/uuid"
	"job-interviewer/internal/interviewer/contracts"
	"job-interviewer/internal/interviewer/model"
	"job-interviewer/internal/interviewer/storage/messages"
	"job-interviewer/pkg/subscription"
)

type freeNextQuestionChecker struct {
	messagesStorage messages.Storage
	userID          uuid.UUID
}

func (f *freeNextQuestionChecker) Check(ctx context.Context) (*subscription.IsAvailableOut, error) {
	messagesFromActiveInterview, err := f.messagesStorage.GetMessagesFromActiveInterviewByUserID(ctx, f.userID)
	if err != nil {
		return nil, err
	}

	questionsCount := 0
	for _, message := range messagesFromActiveInterview {
		if message.Role == model.RoleUser {
			continue
		}

		questionsCount++
	}

	return &subscription.IsAvailableOut{
		Result: questionsCount < defaultFreeQuestions,
		Reason: contracts.ErrQuestionsInFreePlanHaveExpired,
	}, nil
}

func (f *freeNextQuestionChecker) Type() subscription.UserType {
	return subscription.UserTypeFree
}
