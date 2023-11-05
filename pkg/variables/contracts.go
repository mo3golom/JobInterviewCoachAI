package variables

import (
	"time"
)

type (
	Repository interface {
		GetEnvironment() EnvironmentValue
		GetBool(t BoolVariable) bool
		GetInt64(t StringVariable) int64
		GetInt64s(t StringVariable) []int64
		GetString(t StringVariable) string
		GetStrings(t StringVariable) []string
		GetDuration(t StringVariable) time.Duration
		GetConvertedDuration(t StringVariable) (time.Duration, error)
	}
)
