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
			err := g.handleUpdate(ctx, &in, senderImpl)
			if err != nil {
				log.Fatal(err)
			}
		}(update)
	}
}

func (g *DefaultGateway) handleUpdate(ctx context.Context, in *tgbotapi.Update, sender Sender) error {
	request := model.NewRequest(in)

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
		return g.unknownCommand(in)
	}

	return nil
}

func (g *DefaultGateway) determineCommand(request *model.Request) (Handler, bool) {
	command, ok := g.commandHandlers[request.Command]
	return command, ok
}

func (g *DefaultGateway) unknownCommand(in *tgbotapi.Update) error {
	var chatID int64 = 0
	if in.CallbackQuery != nil && in.CallbackQuery.Message != nil && in.CallbackQuery.Message.Chat != nil {
		chatID = in.CallbackQuery.Message.Chat.ID
	}
	if in.Message != nil && in.Message.Chat != nil {
		chatID = in.Message.Chat.ID
	}

	message := tgbotapi.NewMessage(chatID, "Неизвестная команда")
	_, err := g.client.Send(message)
	return err
}
