package model

import "github.com/google/uuid"

const (
	InterviewQuestionStatusAnswered InterviewQuestionStatus = "answered"
	InterviewQuestionStatusActive   InterviewQuestionStatus = "active"
	InterviewQuestionStatusCanceled InterviewQuestionStatus = "canceled"
	InterviewQuestionStatusBad      InterviewQuestionStatus = "bad"
	InterviewQuestionStatusSkip     InterviewQuestionStatus = "skip"
)

type (
	InterviewQuestionStatus string

	Question struct {
		ID      uuid.UUID
		Text    string
		JobInfo JobInfo
	}
)
