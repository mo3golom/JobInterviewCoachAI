package getinterviewoptions

import (
	"job-interviewer/internal/interviewer/contracts"
)

var (
	positions = map[string]string{
		"0": "golang developer",
		"1": "python developer",
		"2": "php developer",
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
	}
}
