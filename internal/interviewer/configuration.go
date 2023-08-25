package interviewer

import (
	"github.com/jmoiron/sqlx"
	"job-interviewer/internal/interviewer/contracts"
	"job-interviewer/internal/interviewer/flow"
	"job-interviewer/internal/interviewer/gpt"
	interview2 "job-interviewer/internal/interviewer/service/interview"
	"job-interviewer/internal/interviewer/storage/interview"
	"job-interviewer/internal/interviewer/storage/messages"
	"job-interviewer/internal/interviewer/storage/user"
	"job-interviewer/internal/interviewer/usecase/acceptanswer"
	"job-interviewer/internal/interviewer/usecase/finishinterview"
	"job-interviewer/internal/interviewer/usecase/getinterview"
	"job-interviewer/internal/interviewer/usecase/getnextquestion"
	"job-interviewer/internal/interviewer/usecase/startinterview"
	user2 "job-interviewer/internal/interviewer/usecase/user"
	"job-interviewer/pkg/transactional"
)

type (
	ConfigurationUseCases struct {
		StartInterview  contracts.StartInterviewUseCase
		FinishInterview contracts.FinishInterviewUseCase
		GetNextQuestion contracts.GetNextQuestionUseCase
		AcceptAnswer    contracts.AcceptAnswerUseCase
		GetInterview    contracts.GetInterviewUsecase
		User            contracts.UserUseCase
	}

	Configuration struct {
		UseCases *ConfigurationUseCases
	}
)

func NewConfiguration(db *sqlx.DB, transactionalTemplate transactional.Template, gptGateway gpt.Gateway) *Configuration {
	interviewStorage := interview.NewStorage(db)
	userStorage := user.NewStorage(db)
	messagesStorage := messages.NewStorage(db)

	interviewService := interview2.NewInterviewService(gptGateway, interviewStorage, messagesStorage, transactionalTemplate)
	interviewFlow := flow.NewDefaultInterviewFlow(interviewService)

	useCases := &ConfigurationUseCases{
		StartInterview:  startinterview.NewUseCase(interviewService, interviewFlow),
		FinishInterview: finishinterview.NewUseCase(interviewFlow),
		GetNextQuestion: getnextquestion.NewUseCase(interviewFlow),
		AcceptAnswer:    acceptanswer.NewUseCase(interviewFlow),
		GetInterview:    getinterview.NewUseCase(interviewService),
		User:            user2.NewUseCase(userStorage, transactionalTemplate),
	}

	return &Configuration{UseCases: useCases}
}
