package model

const (
	BehavioralPosition Position = "behavioral"
)

type (
	Position string

	JobInfo struct {
		Position Position
	}
)
