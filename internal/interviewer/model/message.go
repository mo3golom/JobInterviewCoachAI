package model

import "time"

const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

type (
	Role string

	Message struct {
		Role      Role
		Content   string
		CreatedAt time.Time
	}
)
