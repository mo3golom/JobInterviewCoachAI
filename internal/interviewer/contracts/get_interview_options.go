package contracts

import (
	"job-interviewer/internal/interviewer/model"
)

type GetInterviewOptionsOut struct {
	Positions map[string]string
	Levels    map[string]model.JobLevel
}

type GetInterviewOptionsUseCase interface {
	GetInterviewOptions() GetInterviewOptionsOut
}
