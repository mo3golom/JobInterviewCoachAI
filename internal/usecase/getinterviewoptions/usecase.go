package getinterviewoptions

import (
	"job-interviewer/internal/contracts"
	"job-interviewer/internal/model"
)

var (
	positions = []string{
		"golang developer",
		"python developer",
		"php developer",
	}

	levels = []model.JobLevel{
		model.JobLevelJunior,
		model.JobLevelMiddle,
		model.JobLevelSenior,
	}
)

type UseCase struct {
}

func NewUseCase() *UseCase {
	return &UseCase{}
}

func (u *UseCase) GetInterviewOptions() contracts.GetInterviewOptionsOut {
	return contracts.GetInterviewOptionsOut{
		Positions: positions,
		Levels:    levels,
	}
}
