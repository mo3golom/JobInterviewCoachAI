package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
	jobinterviewer "job-interviewer"
	"job-interviewer/cmd"
	"job-interviewer/internal/interviewer"
	"job-interviewer/internal/interviewer/gpt"
	"job-interviewer/internal/telegram"
	"job-interviewer/pkg/subscription"
	telegramPkg "job-interviewer/pkg/telegram"
	"job-interviewer/pkg/transactional"
	variables2 "job-interviewer/pkg/variables"
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

	log := cmd.MustInitLogger()
	variables, err := variables2.NewConfiguration()
	if err != nil {
		panic(err)
	}

	c := openai.NewClient(
		variables.Repository.MustGet().GetString(jobinterviewer.GPTApiKey),
	)
	gptGateway := gpt.NewGateway(c)

	tgPkgConfig := telegramPkg.NewConfiguration(
		log,
		variables.Repository.MustGet().GetString(jobinterviewer.TGBotToken),
	)
	tgPkg := tgPkgConfig.Gateway

	subscriptionService := subscription.NewSubscriptionService(db)

	interviewerConfig := interviewer.NewConfiguration(
		db,
		template,
		gptGateway,
		subscriptionService,
	)
	telegramConfig := telegram.NewConfiguration(
		interviewerConfig,
		log,
	)

	// REGISTER MIDDLEWARE
	tgPkg.RegisterMiddleware(telegramConfig.Middlewares.User)

	// REGISTER COMMAND
	tgPkg.RegisterCommandHandler(telegramConfig.Handlers.Start)
	tgPkg.RegisterCommandHandler(telegramConfig.Handlers.StartInterview)
	tgPkg.RegisterCommandHandler(telegramConfig.Handlers.FinishInterview)
	tgPkg.RegisterCommandHandler(telegramConfig.Handlers.GetNextQuestion)
	tgPkg.RegisterCommandHandler(telegramConfig.Handlers.SkipQuestion)
	tgPkg.RegisterCommandHandler(telegramConfig.Handlers.GetAnswerSuggestion)

	// REGISTER MESSAGE HANDLER
	tgPkg.RegisterHandler(telegramConfig.Handlers.AcceptAnswer)

	// REGISTER ERROR HANDLER
	tgPkg.RegisterErrorHandler(telegramConfig.Handlers.SubscribeErrorHandler)

	updateConfig := telegramPkg.Config{
		Timeout: 60,
	}
	tgPkg.Run(ctx, updateConfig)
}
