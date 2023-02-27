package internal

import (
	"github.com/jmoiron/sqlx"
	"job-interviewer/internal/contracts"
	"job-interviewer/internal/gpt"
	"job-interviewer/internal/middleware/tguser"
	interview2 "job-interviewer/internal/service/interview"
	question2 "job-interviewer/internal/service/question"
	"job-interviewer/internal/storage/interview"
	"job-interviewer/internal/storage/question"
	"job-interviewer/internal/storage/user"
	"job-interviewer/internal/usecase/acceptanswer"
	"job-interviewer/internal/usecase/finishinterview"
	"job-interviewer/internal/usecase/getinterviewoptions"
	"job-interviewer/internal/usecase/getnextquestion"
	"job-interviewer/internal/usecase/startinterview"
	"job-interviewer/pkg/telegram"
	"job-interviewer/pkg/transactional"
)

type ConfigurationUseCases struct {
	StartInterview      contracts.StartInterviewUsecase
	FinishInterview     contracts.FinishInterviewUseCase
	GetNextQuestion     contracts.GetNextQuestionUseCase
	AcceptAnswer        contracts.AcceptAnswerUseCase
	GetInterviewOptions contracts.GetInterviewOptionsUseCase
}

type Middlewares struct {
	TgUser telegram.Middleware
}

type Configuration struct {
	UseCases    *ConfigurationUseCases
	Middlewares *Middlewares
}

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
	}
	middlewares := &Middlewares{
		TgUser: tguser.NewMiddleware(userStorage, transactionalTemplate),
	}

	return &Configuration{UseCases: useCases, Middlewares: middlewares}
}
