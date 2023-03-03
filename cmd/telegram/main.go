package main

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	gogpt "github.com/sashabaranov/go-gpt3"
	"job-interviewer/cmd"
	"job-interviewer/internal"
	"job-interviewer/internal/gpt"
	"job-interviewer/pkg/telegram"
	"job-interviewer/pkg/telegram/handlers/command/finishinterview"
	"job-interviewer/pkg/telegram/handlers/command/getnextquestion"
	"job-interviewer/pkg/telegram/handlers/command/prestartinterview"
	"job-interviewer/pkg/telegram/handlers/command/start"
	"job-interviewer/pkg/telegram/handlers/command/startinterview"
	"job-interviewer/pkg/telegram/handlers/message/acceptanswer"
	"job-interviewer/pkg/transactional"
	"os"
)

func main() {
	ctx := context.Background()
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	db := cmd.MustInitDB(ctx)
	template := transactional.NewTemplate(db)

	c := gogpt.NewClient(os.Getenv("GPT_API_KEY"))
	gptGateway := gpt.NewGateway(c)

	tgConfig := telegram.NewConfiguration()
	tg := tgConfig.Gateway

	interviewConfig := internal.NewConfiguration(
		db,
		template,
		gptGateway,
	)

	// HANDLERS
	finishInterviewHandler := finishinterview.NewHandler(interviewConfig.UseCases.FinishInterview)
	acceptAnswerHandler := acceptanswer.NewHandler(
		interviewConfig.UseCases.GetNextQuestion,
		interviewConfig.UseCases.AcceptAnswer,
		finishInterviewHandler,
		tgConfig.KeyboardService,
	)
	getNextQuestionHandler := getnextquestion.NewHandler(
		interviewConfig.UseCases.GetNextQuestion,
		finishInterviewHandler,
	)
	startInterviewHandler := startinterview.NewHandler(
		tgConfig.KeyboardService,
		interviewConfig.UseCases.GetInterviewOptions,
		interviewConfig.UseCases.StartInterview,
		getNextQuestionHandler,
	)
	preStartInterviewHandler := prestartinterview.NewHandler(
		tgConfig.KeyboardService,
		interviewConfig.UseCases.GetInterview,
		startInterviewHandler,
		getNextQuestionHandler,
	)

	// REGISTER MIDDLEWARE
	tg.RegisterMiddleware(interviewConfig.Middlewares.TgUser)

	// REGISTER COMMAND
	tg.RegisterCommandHandler(start.NewHandler(tgConfig.KeyboardService))
	tg.RegisterCommandHandler(startInterviewHandler)
	tg.RegisterCommandHandler(preStartInterviewHandler)
	tg.RegisterCommandHandler(finishInterviewHandler)
	tg.RegisterCommandHandler(getNextQuestionHandler)

	// REGISTER MESSAGE HANDLER
	tg.RegisterHandler(acceptAnswerHandler)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	tg.Run(ctx, updateConfig)
}
