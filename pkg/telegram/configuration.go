package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"job-interviewer/pkg/telegram/service/keyboard"
	"os"
)

type Configuration struct {
	KeyboardService keyboard.Service
	Gateway         Gateway
}

func NewConfiguration() *Configuration {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TG_BOT_TOKEN"))
	if err != nil {
		panic(err)
	}

	if os.Getenv("ENV") != "prod" {
		bot.Debug = true
	}

	return &Configuration{
		KeyboardService: &keyboard.DefaultService{},
		Gateway:         NewGateway(bot),
	}
}
