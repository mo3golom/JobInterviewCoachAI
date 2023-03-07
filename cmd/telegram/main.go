package main

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	gogpt "github.com/sashabaranov/go-gpt3"
	"job-interviewer/cmd"
	"job-interviewer/internal/interviewer"
	"job-interviewer/internal/interviewer/gpt"
	"job-interviewer/internal/telegram"
	telegramPkg "job-interviewer/pkg/telegram"
	"job-interviewer/pkg/transactional"
	"os"
)

func main() {
	ctx := context.Background()
	if _, err := os.Stat(".env"); err == nil {
		// path/to/whatever exists
		err := godotenv.Load(".env")
		if err != nil {
			panic(err)
		}
	}

	db := cmd.MustInitDB(ctx)
	template := transactional.NewTemplate(db)

	c := gogpt.NewClient(os.Getenv("GPT_API_KEY"))
	gptGateway := gpt.NewGateway(c)

	tgPkgConfig := telegramPkg.NewConfiguration()
	tgPkg := tgPkgConfig.Gateway

	interviewerConfig := interviewer.NewConfiguration(
		db,
		template,
		gptGateway,
	)
	telegramConfig := telegram.NewConfiguration(
		interviewerConfig,
		tgPkgConfig,
		db,
	)

	// REGISTER MIDDLEWARE
	tgPkg.RegisterMiddleware(telegramConfig.Middlewares.User)

	// REGISTER COMMAND
	tgPkg.RegisterCommandHandler(telegramConfig.Handlers.Start)
	tgPkg.RegisterCommandHandler(telegramConfig.Handlers.StartInterview)
	tgPkg.RegisterCommandHandler(telegramConfig.Handlers.PreStartInterview)
	tgPkg.RegisterCommandHandler(telegramConfig.Handlers.FinishInterview)
	tgPkg.RegisterCommandHandler(telegramConfig.Handlers.GetNextQuestion)
	tgPkg.RegisterCommandHandler(telegramConfig.Handlers.MarkQuestionAsBad)
	tgPkg.RegisterCommandHandler(telegramConfig.Handlers.MarkQuestionAsSkip)

	// REGISTER MESSAGE HANDLER
	tgPkg.RegisterHandler(telegramConfig.Handlers.AcceptAnswer)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	tgPkg.Run(ctx, updateConfig)
}
