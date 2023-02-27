package contracts

import "job-interviewer/internal/model"

type GetInterviewOptionsOut struct {
	Positions []string
	Levels    []model.JobLevel
}

type GetInterviewOptionsUseCase interface {
	GetInterviewOptions() GetInterviewOptionsOut
}
