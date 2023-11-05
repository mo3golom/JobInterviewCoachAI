package model

import "github.com/google/uuid"

const (
	InterviewStatusCreated  InterviewStatus = "created"
	InterviewStatusStarted  InterviewStatus = "started"
	InterviewStatusFinished InterviewStatus = "finished"

	InterviewStateStart           InterviewState = "start"
	InterviewStateWaitingQuestion InterviewState = "waiting_question"
	InterviewStateWaitingAnswer   InterviewState = "waiting_answer"
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
