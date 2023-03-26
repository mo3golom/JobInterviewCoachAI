package telegram

import (
	"github.com/jmoiron/sqlx"
	"job-interviewer/internal/interviewer"
	"job-interviewer/internal/telegram/handlers/command/changeuserlanguage"
	"job-interviewer/internal/telegram/handlers/command/finishinterview"
	"job-interviewer/internal/telegram/handlers/command/getnextquestion"
	"job-interviewer/internal/telegram/handlers/command/prestartinterview"
	"job-interviewer/internal/telegram/handlers/command/questionbad"
	"job-interviewer/internal/telegram/handlers/command/questionskip"
	"job-interviewer/internal/telegram/handlers/command/start"
	"job-interviewer/internal/telegram/handlers/command/startinterview"
	"job-interviewer/internal/telegram/handlers/message/acceptanswer"
	"job-interviewer/internal/telegram/language"
	"job-interviewer/internal/telegram/language/en"
	"job-interviewer/internal/telegram/language/ru"
	"job-interviewer/internal/telegram/middleware/user"
	tgService "job-interviewer/internal/telegram/service"
	storage2 "job-interviewer/internal/telegram/storage"
	language2 "job-interviewer/pkg/language"
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
		ChangeUserLanguage telegram.CommandHandler

		AcceptAnswer telegram.Handler
	}

	Middlewares struct {
		User telegram.Middleware
	}

	Configuration struct {
		Handlers        *ConfigurationHandlers
		Middlewares     *Middlewares
		LanguageService language.Service
	}
)

func NewConfiguration(
	interviewerConfig *interviewer.Configuration,
	tgConfig *telegram.Configuration,
	db *sqlx.DB,
) *Configuration {
	storage := storage2.NewStorage(db)
	languageService := language.NewService(map[language2.Language]language.Dictionary{
		language.English: en.Dict{},
		language.Russian: ru.Dict{},
	})

	service := tgService.NewService(
		interviewerConfig.UseCases.FinishInterview,
		interviewerConfig.UseCases.GetNextQuestion,
		tgConfig.KeyboardService,
		storage,
		languageService,
	)

	startInterviewHandler := startinterview.NewHandler(
		tgConfig.KeyboardService,
		interviewerConfig.UseCases.GetInterviewOptions,
		interviewerConfig.UseCases.StartInterview,
		service,
		languageService,
	)
	preStartInterviewHandler := prestartinterview.NewHandler(
		tgConfig.KeyboardService,
		interviewerConfig.UseCases.GetInterview,
		startInterviewHandler,
		service,
		languageService,
	)

	configurationHandlers := &ConfigurationHandlers{
		Start:              start.NewHandler(service),
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
			languageService,
		),
		ChangeUserLanguage: changeuserlanguage.NewHandler(
			interviewerConfig.UseCases.User,
			service,
			tgConfig.KeyboardService,
			languageService,
		),
	}

	middlewares := &Middlewares{
		User: user.NewMiddleware(interviewerConfig.UseCases.User),
	}

	return &Configuration{
		Handlers:        configurationHandlers,
		Middlewares:     middlewares,
		LanguageService: languageService,
	}
}
