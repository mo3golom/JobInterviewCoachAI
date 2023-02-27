package model

import "github.com/google/uuid"

const (
	InterviewStatusStarted  InterviewStatus = "started"
	InterviewStatusFinished InterviewStatus = "finished"
)

type (
	InterviewStatus string

	Interview struct {
		ID             uuid.UUID
		UserID         uuid.UUID
		Status         InterviewStatus
		JobInfo        JobInfo
		QuestionsCount int64
	}
)
