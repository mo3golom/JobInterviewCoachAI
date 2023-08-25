package telegram

import (
	"job-interviewer/internal/interviewer"
	"job-interviewer/internal/telegram/handlers/command/finishinterview"
	"job-interviewer/internal/telegram/handlers/command/getnextquestion"
	"job-interviewer/internal/telegram/handlers/command/start"
	"job-interviewer/internal/telegram/handlers/command/startinterview"
	"job-interviewer/internal/telegram/handlers/message/acceptanswer"
	"job-interviewer/internal/telegram/middleware/user"
	tgService "job-interviewer/internal/telegram/service"
	"job-interviewer/pkg/telegram"
)

type (
	ConfigurationHandlers struct {
		Start           telegram.CommandHandler
		StartInterview  telegram.CommandHandler
		FinishInterview telegram.CommandHandler
		GetNextQuestion telegram.CommandHandler

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
) *Configuration {
	service := tgService.NewService(
		interviewerConfig.UseCases.FinishInterview,
		interviewerConfig.UseCases.GetNextQuestion,
	)

	startInterviewHandler := startinterview.NewHandler(
		interviewerConfig.UseCases.StartInterview,
		interviewerConfig.UseCases.GetInterview,
		service,
	)

	configurationHandlers := &ConfigurationHandlers{
		Start:           start.NewHandler(service),
		StartInterview:  startInterviewHandler,
		FinishInterview: finishinterview.NewHandler(service),
		GetNextQuestion: getnextquestion.NewHandler(service),
		AcceptAnswer: acceptanswer.NewHandler(
			interviewerConfig.UseCases.AcceptAnswer,
			service,
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
