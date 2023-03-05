package getinterviewoptions

import (
	"job-interviewer/internal/interviewer/contracts"
	"job-interviewer/internal/interviewer/model"
)

var (
	positions = map[string]string{
		"0": "golang developer",
		"1": "python developer",
		"2": "php developer",
	}

	levels = map[string]model.JobLevel{
		"0": model.JobLevelJunior,
		"1": model.JobLevelMiddle,
		"2": model.JobLevelSenior,
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
