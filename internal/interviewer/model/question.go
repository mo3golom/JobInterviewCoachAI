package model

import "github.com/google/uuid"

type (
	InterviewQuestionStatus string

	Question struct {
		ID      uuid.UUID
		Text    string
		JobInfo JobInfo
	}

	InterviewQuestion struct {
		Text   string
		Answer string
	}

	AnswerSuggestion struct {
		Text string
	}
)
