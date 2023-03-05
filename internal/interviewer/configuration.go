package interviewer

import (
	"github.com/jmoiron/sqlx"
	"job-interviewer/internal/interviewer/contracts"
	"job-interviewer/internal/interviewer/gpt"
	interview2 "job-interviewer/internal/interviewer/service/interview"
	question2 "job-interviewer/internal/interviewer/service/question"
	"job-interviewer/internal/interviewer/storage/interview"
	"job-interviewer/internal/interviewer/storage/question"
	"job-interviewer/internal/interviewer/storage/user"
	"job-interviewer/internal/interviewer/usecase/acceptanswer"
	"job-interviewer/internal/interviewer/usecase/finishinterview"
	"job-interviewer/internal/interviewer/usecase/getinterview"
	"job-interviewer/internal/interviewer/usecase/getinterviewoptions"
	"job-interviewer/internal/interviewer/usecase/getnextquestion"
	"job-interviewer/internal/interviewer/usecase/startinterview"
	"job-interviewer/internal/interviewer/usecase/updatequestion"
	user2 "job-interviewer/internal/interviewer/usecase/user"
	"job-interviewer/pkg/transactional"
)

type (
	ConfigurationUseCases struct {
		StartInterview      contracts.StartInterviewUseCase
		FinishInterview     contracts.FinishInterviewUseCase
		GetNextQuestion     contracts.GetNextQuestionUseCase
		AcceptAnswer        contracts.AcceptAnswerUseCase
		GetInterviewOptions contracts.GetInterviewOptionsUseCase
		GetInterview        contracts.GetInterviewUseCase
		UpdateQuestion      contracts.UpdateQuestionUseCase
		User                contracts.UserUseCase
	}

	Configuration struct {
		UseCases *ConfigurationUseCases
	}
)

func NewConfiguration(db *sqlx.DB, transactionalTemplate transactional.Template, gptGateway gpt.Gateway) *Configuration {
	questionStorage := question.NewStorage(db)
	interviewStorage := interview.NewStorage(db)
	userStorage := user.NewStorage(db)

	questionService := question2.NewQuestionService(gptGateway, questionStorage, transactionalTemplate)
	interviewService := interview2.NewInterviewService(gptGateway, interviewStorage, questionStorage, transactionalTemplate)

	useCases := &ConfigurationUseCases{
		StartInterview:      startinterview.NewUseCase(interviewService, questionService),
		FinishInterview:     finishinterview.NewUseCase(interviewService),
		GetNextQuestion:     getnextquestion.NewUseCase(interviewService, questionService),
		AcceptAnswer:        acceptanswer.NewUseCase(interviewService),
		GetInterviewOptions: getinterviewoptions.NewUseCase(),
		GetInterview:        getinterview.NewUseCase(interviewService),
		UpdateQuestion:      updatequestion.NewUseCase(interviewStorage, questionStorage, transactionalTemplate),
		User:                user2.NewUseCase(userStorage, transactionalTemplate),
	}

	return &Configuration{UseCases: useCases}
}
