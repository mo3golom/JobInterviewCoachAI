package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"job-interviewer/pkg/logger"
	"os"
)

type Configuration struct {
	Gateway Gateway
}

func NewConfiguration(log logger.Logger, token string) *Configuration {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}

	if os.Getenv("ENV") != "prod" {
		bot.Debug = true
	}

	return &Configuration{
		Gateway: NewGateway(bot, log),
	}
}
