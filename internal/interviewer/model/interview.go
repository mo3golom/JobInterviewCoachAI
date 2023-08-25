package model

import "github.com/google/uuid"

const (
	InterviewStatusCreated  InterviewStatus = "created"
	InterviewStatusStarted  InterviewStatus = "started"
	InterviewStatusFinished InterviewStatus = "finished"

	InterviewStateDefault             InterviewState = "default"
	InterviewStateWaitingQuestion     InterviewState = "waiting_question"
	InterviewStateAnsweringOnQuestion InterviewState = "answering_on_question"
)

type (
	InterviewStatus string
	InterviewState  string

	Interview struct {
		ID             uuid.UUID
		UserID         uuid.UUID
		Status         InterviewStatus
		State          InterviewState
		JobInfo        JobInfo
		QuestionsCount int64
	}
)
