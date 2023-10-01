package contracts

import "errors"

var (
	ErrEmptyActiveInterview           = errors.New("active interview not found")
	ErrInterviewQuestionsIsEmpty      = errors.New("interview questions is empty")
	ErrActionDoesntAllow              = errors.New("action doesn't allow in this state")
	ErrFreeAttemptsHaveExpired        = errors.New("free attempts have expired")
	ErrQuestionsInFreePlanHaveExpired = errors.New("questions in free plan have expired")
	ErrPaidSubscriptionHasExpired     = errors.New("paid subscription has expired")
)
