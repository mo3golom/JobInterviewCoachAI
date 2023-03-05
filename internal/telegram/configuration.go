package telegram

import (
	"github.com/jmoiron/sqlx"
	"job-interviewer/internal/interviewer"
	"job-interviewer/internal/telegram/handlers/command/finishinterview"
	"job-interviewer/internal/telegram/handlers/command/getnextquestion"
	"job-interviewer/internal/telegram/handlers/command/prestartinterview"
	"job-interviewer/internal/telegram/handlers/command/questionbad"
	"job-interviewer/internal/telegram/handlers/command/questionskip"
	"job-interviewer/internal/telegram/handlers/command/start"
	"job-interviewer/internal/telegram/handlers/command/startinterview"
	"job-interviewer/internal/telegram/handlers/message/acceptanswer"
	"job-interviewer/internal/telegram/middleware/user"
	tgService "job-interviewer/internal/telegram/service"
	storage2 "job-interviewer/internal/telegram/storage"
	"job-interviewer/pkg/telegram"
)

type (
	ConfigurationHandlers struct {
		Start              telegram.CommandHandler
		PreStartInterview  telegram.CommandHandler
		StartInterview     telegram.CommandHandler
		FinishInterview    telegram.CommandHandler
		GetNextQuestion    telegram.CommandHandler
		MarkQuestionAsBad  telegram.CommandHandler
		MarkQuestionAsSkip telegram.CommandHandler

		AcceptAnswer telegram.Handler
	}

	Middlewares struct {
		User telegram.Middleware
	}

	Configuration struct {
		Handlers    *ConfigurationHandlers
		Middlewares *Middlewares
	}
)

func NewConfiguration(
	interviewerConfig *interviewer.Configuration,
	tgConfig *telegram.Configuration,
	db *sqlx.DB,
) *Configuration {
	storage := storage2.NewStorage(db)

	service := tgService.NewService(
		interviewerConfig.UseCases.FinishInterview,
		interviewerConfig.UseCases.GetNextQuestion,
		tgConfig.KeyboardService,
		storage,
	)

	startInterviewHandler := startinterview.NewHandler(
		tgConfig.KeyboardService,
		interviewerConfig.UseCases.GetInterviewOptions,
		interviewerConfig.UseCases.StartInterview,
		service,
	)
	preStartInterviewHandler := prestartinterview.NewHandler(
		tgConfig.KeyboardService,
		interviewerConfig.UseCases.GetInterview,
		startInterviewHandler,
		service,
	)

	configurationHandlers := &ConfigurationHandlers{
		Start:              start.NewHandler(tgConfig.KeyboardService),
		PreStartInterview:  preStartInterviewHandler,
		StartInterview:     startInterviewHandler,
		FinishInterview:    finishinterview.NewHandler(service),
		GetNextQuestion:    getnextquestion.NewHandler(service),
		MarkQuestionAsBad:  questionbad.NewHandler(interviewerConfig.UseCases.UpdateQuestion, service),
		MarkQuestionAsSkip: questionskip.NewHandler(interviewerConfig.UseCases.UpdateQuestion, service),
		AcceptAnswer: acceptanswer.NewHandler(
			interviewerConfig.UseCases.AcceptAnswer,
			service,
			tgConfig.KeyboardService,
		),
	}

	middlewares := &Middlewares{
		User: user.NewMiddleware(interviewerConfig.UseCases.User),
	}

	return &Configuration{
		Handlers:    configurationHandlers,
		Middlewares: middlewares,
	}
}
