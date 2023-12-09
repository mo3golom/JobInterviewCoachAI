package main

import (
	"context"
	"github.com/joho/godotenv"
	jobinterviewer "job-interviewer"
	"job-interviewer/cmd"
	"job-interviewer/internal/interviewer"
	"job-interviewer/internal/interviewer/gpt"
	"job-interviewer/internal/telegram"
	go_openai "job-interviewer/pkg/go-openai"
	"job-interviewer/pkg/payments"
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

	c := go_openai.NewClient(go_openai.Config{
		ServiceType: go_openai.ServiceType(variables.Repository.MustGet().GetString(jobinterviewer.GPTServiceType)),
		OpenAI: &go_openai.OpenAIConfig{
			AuthToken: variables.Repository.MustGet().GetString(jobinterviewer.GPTApiKey),
		},
		VseGPT: &go_openai.VseGPTConfig{
			AuthToken: variables.Repository.MustGet().GetString(jobinterviewer.VseGPTApiKey),
			BaseUrl:   variables.Repository.MustGet().GetString(jobinterviewer.VseGPTBaseUrl),
		},
	})
	gptGateway := gpt.NewGateway(c)

	tgPkgConfig := telegramPkg.NewConfiguration(
		log,
		variables.Repository.MustGet().GetString(jobinterviewer.TGBotToken),
	)
	tgPkg := tgPkgConfig.Gateway

	subscriptionService := subscription.NewSubscriptionService(
		db,
		variables.Repository.MustGet().GetInt64(jobinterviewer.FreeInterviewsCount),
	)
	paymentsService := payments.NewPaymentsService(
		db,
		template,
		variables.Repository.MustGet().GetInt64(jobinterviewer.YMShopID),
		variables.Repository.MustGet().GetString(jobinterviewer.YMSecretKey),
	)

	interviewerConfig := interviewer.NewConfiguration(
		db,
		template,
		gptGateway,
		subscriptionService,
		paymentsService,
		variables.Repository.MustGet(),
	)
	telegramConfig := telegram.NewConfiguration(
		interviewerConfig,
		log,
		paymentsService,
		variables.Repository.MustGet(),
	)

	// REGISTER MIDDLEWARE
	tgPkg.RegisterMiddleware(telegramConfig.Middlewares.User)

	// REGISTER COMMAND
	tgPkg.RegisterCommandHandler(
		telegramConfig.Handlers.Start,
		telegramConfig.Handlers.StartInterview,
		telegramConfig.Handlers.FinishInterview,
		telegramConfig.Handlers.GetNextQuestion,
		telegramConfig.Handlers.SkipQuestion,
		telegramConfig.Handlers.GetAnswerSuggestion,
		telegramConfig.Handlers.PaySubscription,
		telegramConfig.Handlers.CheckPayment,
		telegramConfig.Handlers.About,
	)

	// REGISTER MESSAGE HANDLER
	tgPkg.RegisterHandler(telegramConfig.Handlers.AcceptAnswer)

	// REGISTER ERROR HANDLER
	tgPkg.RegisterErrorHandler(telegramConfig.Handlers.SubscribeErrorHandler)

	updateConfig := telegramPkg.Config{
		Timeout: 60,
	}
	tgPkg.Run(ctx, updateConfig)
}
