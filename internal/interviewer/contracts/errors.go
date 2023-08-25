package contracts

import "errors"

var (
	ErrEmptyActiveInterview      = errors.New("active interview not found")
	ErrInterviewQuestionsIsEmpty = errors.New("interview questions is empty")
)
