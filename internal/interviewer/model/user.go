package model

import (
	"github.com/google/uuid"
	"job-interviewer/pkg/language"
)

type (
	User struct {
		ID   uuid.UUID
		Lang language.Language
	}
)
