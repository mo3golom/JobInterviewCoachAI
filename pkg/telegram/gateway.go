package telegram

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"job-interviewer/pkg/helper"
	"job-interviewer/pkg/telegram/model"
	"log"
)

type DefaultGateway struct {
	client          externalClient
	commandHandlers map[string]Handler
	messageHandlers []Handler
	middlewares     []Middleware
}

func NewGateway(c externalClient) *DefaultGateway {
	return &DefaultGateway{client: c}
}

func (g *DefaultGateway) RegisterMiddleware(middleware Middleware) {
	copyMiddleware := helper.CopySlice[Middleware](g.middlewares)
	copyMiddleware = append(copyMiddleware, middleware)

	g.middlewares = copyMiddleware
}

func (g *DefaultGateway) RegisterCommandHandler(handler CommandHandler) {
	copyCommandHandlers := helper.CopyMap[string, Handler](g.commandHandlers)

	copyCommandHandlers[handler.Command()] = handler
	g.commandHandlers = copyCommandHandlers
}

func (g *DefaultGateway) RegisterHandler(handler Handler) {
	copyMessageHandlers := helper.CopySlice[Handler](g.messageHandlers)
	copyMessageHandlers = append(copyMessageHandlers, handler)

	g.messageHandlers = copyMessageHandlers
}

func (g *DefaultGateway) Run(ctx context.Context, config tgbotapi.UpdateConfig) {
	senderImpl := &sender{
		client: g.client,
	}

	updates := g.client.GetUpdatesChan(config)
	for update := range updates {
		go func(in tgbotapi.Update) {
			request := model.NewRequest(in)
			err := g.handleUpdate(ctx, &request, senderImpl)
			if err != nil {
				log.Println("ERR: ", err)

				_, _ = senderImpl.Send(
					model.
						NewResponse(request.Chat.ID).
						SetText("Что-то пошло не так :("),
				)
			}
		}(update)
	}
}

func (g *DefaultGateway) handleUpdate(ctx context.Context, request *model.Request, sender Sender) error {
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
		return command.Handle(ctx, request, sender)
	}

	// MESSAGE HANDLER
	for _, handler := range g.messageHandlers {
		err := handler.Handle(ctx, request, sender)
		if err == nil {
			continue
		}

		return err
	}

	if len(g.messageHandlers) == 0 {
		_, err := sender.Send(
			model.
				NewResponse(request.Chat.ID).
				SetText("Неизвестная команда"),
		)
		return err
	}

	return nil
}

func (g *DefaultGateway) determineCommand(request *model.Request) (Handler, bool) {
	command, ok := g.commandHandlers[request.Command]
	return command, ok
}
