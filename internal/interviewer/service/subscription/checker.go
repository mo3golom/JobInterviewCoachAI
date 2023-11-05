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
	messagesStorage    messages.Storage
	userID             uuid.UUID
	freeQuestionsCount int64
}

func (f *freeNextQuestionChecker) Check(ctx context.Context) (*subscription.IsAvailableOut, error) {
	messagesFromActiveInterview, err := f.messagesStorage.GetMessagesFromActiveInterviewByUserID(ctx, f.userID)
	if err != nil {
		return nil, err
	}

	var questionsCount int64
	for _, message := range messagesFromActiveInterview {
		if message.Role == model.RoleUser {
			continue
		}

		questionsCount++
	}

	return &subscription.IsAvailableOut{
		Result: questionsCount < f.freeQuestionsCount,
		Reason: contracts.ErrQuestionsInFreePlanHaveExpired,
	}, nil
}

func (f *freeNextQuestionChecker) Type() subscription.UserType {
	return subscription.UserTypeFree
}
