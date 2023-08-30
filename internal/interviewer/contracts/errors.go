package contracts

import "errors"

var (
	ErrEmptyActiveInterview      = errors.New("active interview not found")
	ErrInterviewQuestionsIsEmpty = errors.New("interview questions is empty")
	ErrActionDoesntAllow         = errors.New("action doesn't allow in this state")
)
