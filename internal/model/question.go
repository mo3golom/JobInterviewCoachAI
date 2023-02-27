package model

import "github.com/google/uuid"

const (
	InterviewQuestionStatusCreated  InterviewQuestionStatus = "created"
	InterviewQuestionStatusAnswered InterviewQuestionStatus = "answered"
	InterviewQuestionStatusActive   InterviewQuestionStatus = "active"
)

type (
	InterviewQuestionStatus string

	Question struct {
		ID      uuid.UUID
		Text    string
		JobInfo JobInfo
	}
)
