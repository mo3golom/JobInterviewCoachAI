package telegram

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"job-interviewer/pkg/helper"
	"job-interviewer/pkg/logger"
	"job-interviewer/pkg/telegram/model"
	"net/http"
)

const (
	webhookRoute = "/updates"
)

type DefaultGateway struct {
	client          externalClient
	commandHandlers map[string]Handler
	errorHandlers   []ErrorHandler
	messageHandlers []Handler
	middlewares     []Middleware

	errorLogHandler ErrorHandler
	token           string
}

func NewGateway(
	c externalClient,
	log logger.Logger,
	token string,
) *DefaultGateway {
	gateway := &DefaultGateway{
		client: c,
		token:  token,
	}
	gateway.errorLogHandler = &ErrorLogHandler{log: log}

	return gateway
}

func (g *DefaultGateway) RegisterMiddleware(middleware Middleware) {
	copyMiddleware := helper.CopySlice[Middleware](g.middlewares)
	copyMiddleware = append(copyMiddleware, middleware)

	g.middlewares = copyMiddleware
}

func (g *DefaultGateway) RegisterCommandHandler(handler CommandHandler) {
	copyCommandHandlers := helper.CopyMap[string, Handler](g.commandHandlers)

	copyCommandHandlers[handler.Command()] = handler
	for _, alias := range handler.Aliases() {
		copyCommandHandlers[alias] = handler
	}
	g.commandHandlers = copyCommandHandlers
}

func (g *DefaultGateway) RegisterErrorHandler(handler ErrorHandler) {
	copyErrorHandlers := helper.CopySlice[ErrorHandler](g.errorHandlers)
	copyErrorHandlers = append(copyErrorHandlers, handler)

	g.errorHandlers = copyErrorHandlers
}

func (g *DefaultGateway) RegisterHandler(handler Handler) {
	copyMessageHandlers := helper.CopySlice[Handler](g.messageHandlers)
	copyMessageHandlers = append(copyMessageHandlers, handler)

	g.messageHandlers = copyMessageHandlers
}

func (g *DefaultGateway) Run(ctx context.Context, config Config) {
	senderImpl := &sender{
		client: g.client,
	}

	tgConfig := tgbotapi.NewUpdate(config.Offset)
	tgConfig.Timeout = config.Timeout
	tgConfig.Limit = config.Limit

	var updates tgbotapi.UpdatesChannel
	if config.Webhook != nil && config.Webhook.Enable {
		url := fmt.Sprintf("%s/%s", config.Webhook.Host, g.token)
		if !config.Webhook.Debug {
			webhook, err := tgbotapi.NewWebhook(fmt.Sprintf("%s%s", url, webhookRoute))
			if err != nil {
				panic(err)
			}

			_, err = g.client.Request(webhook)
			if err != nil {
				panic(err)
			}

			info, err := g.client.GetWebhookInfo()
			if err != nil {
				panic(err)
			}

			if info.LastErrorDate != 0 {
				panic(fmt.Sprintf("failed to set webhook: %s", info.LastErrorMessage))
			}
		}

		updates = g.client.ListenForWebhook(webhookRoute)
		go func(addr string) {
			err := http.ListenAndServe(addr, nil)
			if err != nil {
				panic(err)
			}
		}(url)

		fmt.Printf("bot started, webhook: %s", url)
	} else {
		updates = g.client.GetUpdatesChan(tgConfig)
	}

	for update := range updates {
		go func(in tgbotapi.Update) {
			request := model.NewRequest(in)
			err := g.handleUpdate(ctx, &request, senderImpl)
			if err == nil {
				return
			}

			g.handleError(ctx, err, &request, senderImpl)
		}(update)
	}

}
func (g *DefaultGateway) handleUpdate(ctx context.Context, request *model.Request, senderImpl *sender) error {
	senderAdapterImpl := senderAdapter{
		senderImpl,
		request.Chat.ID,
	}
	defer (func() {
		senderAdapterImpl = senderAdapter{}
	})()

	// MIDDLEWARE
	for _, middleware := range g.middlewares {
		err := middleware.Handle(ctx, request)
		if err != nil {
			return err
		}
	}

	// COMMAND
	command, ok := g.determineCommand(request)
	if ok {
		return command.Handle(ctx, request, senderAdapterImpl)
	}

	// MESSAGE HANDLER
	for _, handler := range g.messageHandlers {
		err := handler.Handle(ctx, request, senderAdapterImpl)
		if err == nil {
			continue
		}

		return err
	}

	if len(g.messageHandlers) == 0 {
		_, err := senderAdapterImpl.Send(
			model.
				NewResponse().
				SetText("Unknown command"),
		)
		return err
	}

	return nil
}

func (g *DefaultGateway) handleError(ctx context.Context, err error, request *model.Request, senderImpl *sender) {
	senderAdapterImpl := senderAdapter{
		senderImpl,
		request.Chat.ID,
	}
	defer (func() {
		senderAdapterImpl = senderAdapter{}
	})()

	var ok bool
	for _, errHandler := range g.errorHandlers {
		if errHandler == nil {
			continue
		}

		ok = errHandler.Handle(ctx, err, request, senderAdapterImpl)
		if !ok {
			continue
		}

		return
	}

	if !ok {
		g.errorLogHandler.Handle(ctx, err, request, senderAdapterImpl)
	}
}

func (g *DefaultGateway) determineCommand(request *model.Request) (Handler, bool) {
	command, ok := g.commandHandlers[request.Command]
	return command, ok
}
