package subscription

import (
	"errors"
	"time"
)

const (
	PlanMonth Plan = "month"

	UserTypeFree UserType = "free"
	UserTypePaid UserType = "paid"
)

var (
	ErrFreeAttemptsHaveExpired    = errors.New("free attempts have expired")
	ErrPaidSubscriptionHasExpired = errors.New("paid subscription has expired")
	ErrEmptyUserInfo              = errors.New("empty user info")
)

type (
	Plan     string
	UserType string

	Subscription struct {
		StartAt time.Time
		EndAt   time.Time
	}

	Info struct {
		Type         UserType
		FreeAttempts *int64
		Subscription *Subscription
	}
)

func (u Subscription) IsActive() bool {
	now := time.Now()

	return now.After(u.StartAt) && now.Before(u.EndAt)
}
