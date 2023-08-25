package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"job-interviewer/pkg/logger"
	"os"
)

type Configuration struct {
	Gateway Gateway
}

func NewConfiguration(log logger.Logger) *Configuration {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TG_BOT_TOKEN"))
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
