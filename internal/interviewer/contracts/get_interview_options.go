package contracts

import (
	"job-interviewer/internal/interviewer/model"
)

type GetInterviewOptionsOut struct {
	Positions []string
	Levels    []model.JobLevel
}

type GetInterviewOptionsUseCase interface {
	GetInterviewOptions() GetInterviewOptionsOut
}
