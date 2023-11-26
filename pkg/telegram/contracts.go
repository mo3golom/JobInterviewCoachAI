package telegram

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"job-interviewer/pkg/telegram/model"
)

type (
	externalClient interface {
		GetUpdatesChan(config tgbotapi.UpdateConfig) tgbotapi.UpdatesChannel
		ListenForWebhook(pattern string) tgbotapi.UpdatesChannel
		Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
		Request(c tgbotapi.Chattable) (*tgbotapi.APIResponse, error)
		GetFileDirectURL(fileID string) (string, error)
		GetWebhookInfo() (tgbotapi.WebhookInfo, error)
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
		Aliases() []string
	}

	ErrorHandler interface {
		Handle(ctx context.Context, err error, request *model.Request, sender Sender) bool
	}

	Config struct {
		Offset  int
		Limit   int
		Timeout int
		Webhook *WebhookConfig
	}

	WebhookConfig struct {
		Enable bool
		Host   string
		Debug  bool
	}

	Gateway interface {
		RegisterMiddleware(middleware Middleware)
		RegisterCommandHandler(handler CommandHandler)
		RegisterHandler(handler Handler)
		RegisterErrorHandler(handler ErrorHandler)
		Run(ctx context.Context, config Config)
	}

	Sender interface {
		Send(response model.Response) (int64, error)
		SendCallback(callbackID model.CallbackID, message ...string) error
		Update(messageID int64, response model.Response) error
	}
)
