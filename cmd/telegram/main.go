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
	"job-interviewer/pkg/telegram/handlers/command/start"
	"job-interviewer/pkg/telegram/handlers/command/startinterview"
	"job-interviewer/pkg/telegram/handlers/message/processinterview"
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

	processInterviewHandler := processinterview.NewHandler(
		interviewConfig.UseCases.GetNextQuestion,
		interviewConfig.UseCases.FinishInterview,
		interviewConfig.UseCases.AcceptAnswer,
	)

	// REGISTER MIDDLEWARE
	tg.RegisterMiddleware(interviewConfig.Middlewares.TgUser)

	// REGISTER COMMAND
	tg.RegisterCommandHandler(startinterview.NewHandler(
		tgConfig.KeyboardService,
		interviewConfig.UseCases.GetInterviewOptions,
		interviewConfig.UseCases.StartInterview,
		processInterviewHandler,
	))
	tg.RegisterCommandHandler(start.NewHandler(tgConfig.KeyboardService))

	// REGISTER MESSAGE HANDLER
	tg.RegisterHandler(processInterviewHandler)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	tg.Run(ctx, updateConfig)
}
