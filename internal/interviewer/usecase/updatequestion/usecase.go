package updatequestion

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"job-interviewer/internal/interviewer/model"
	"job-interviewer/internal/interviewer/storage/interview"
	"job-interviewer/internal/interviewer/storage/question"
	"job-interviewer/pkg/transactional"
)

type UseCase struct {
	questionStorage       question.Storage
	interviewStorage      interview.Storage
	transactionalTemplate transactional.Template
}

func NewUseCase(
	is interview.Storage,
	qs question.Storage,
	tr transactional.Template,
) *UseCase {
	return &UseCase{
		interviewStorage:      is,
		questionStorage:       qs,
		transactionalTemplate: tr,
	}
}

func (u *UseCase) MarkActiveQuestionAsBad(ctx context.Context, userID uuid.UUID) error {
	return u.updateStatusForActiveQuestion(ctx, userID, model.InterviewQuestionStatusBad)
}

func (u *UseCase) MarkActiveQuestionAsSkip(ctx context.Context, userID uuid.UUID) error {
	return u.updateStatusForActiveQuestion(ctx, userID, model.InterviewQuestionStatusSkip)
}

func (u *UseCase) updateStatusForActiveQuestion(
	ctx context.Context,
	userID uuid.UUID,
	targetStatus model.InterviewQuestionStatus,
) error {
	return u.transactionalTemplate.Execute(ctx, func(tx transactional.Tx) error {
		activeInterview, err := u.interviewStorage.FindActiveInterviewByUserID(ctx, tx, userID)
		if errors.Is(err, interview.ErrEmptyInterviewResult) {
			return nil
		}
		if err != nil {
			return err
		}

		return u.questionStorage.UpdateInterviewQuestionStatus(
			ctx,
			tx,
			question.UpdateInterviewQuestionStatusIn{
				InterviewID: activeInterview.ID,
				Current:     model.InterviewQuestionStatusActive,
				Target:      targetStatus,
			},
		)
	})
}
