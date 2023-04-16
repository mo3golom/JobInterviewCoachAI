package contracts

type GetInterviewOptionsOut struct {
	Positions map[string]string
}

type GetInterviewOptionsUseCase interface {
	GetInterviewOptions() GetInterviewOptionsOut
}
