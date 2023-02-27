package telegram

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"job-interviewer/pkg/telegram/model"
)

type (
	externalClient interface {
		GetUpdatesChan(config tgbotapi.UpdateConfig) tgbotapi.UpdatesChannel
		Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
		Request(c tgbotapi.Chattable) (*tgbotapi.APIResponse, error)
	}

	Middleware interface {
		Handle(ctx context.Context, request *model.Request) error
	}

	Handler interface {
		Handle(ctx context.Context, request *model.Request, sender Sender) error
	}

	CommandHandler interface {
		Handler
		Command() string
	}

	Gateway interface {
		RegisterMiddleware(middleware Middleware)
		RegisterCommandHandler(handler CommandHandler)
		RegisterHandler(handler Handler)
		Run(ctx context.Context, config tgbotapi.UpdateConfig)
	}

	Sender interface {
		Send(response model.Response) (int64, error)
		Update(messageID int64, response model.Response) error
	}
)
