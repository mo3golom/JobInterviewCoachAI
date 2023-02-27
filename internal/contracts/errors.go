package contracts

import "errors"

var (
	ErrNextQuestionEmpty         = errors.New("next question is empty")
	ErrEmptyActiveInterview      = errors.New("active interview not found")
	ErrInterviewQuestionsIsEmpty = errors.New("interview questions is empty")
)
