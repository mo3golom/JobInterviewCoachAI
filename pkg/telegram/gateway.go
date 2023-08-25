package telegram

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"job-interviewer/pkg/helper"
	"job-interviewer/pkg/logger"
	"job-interviewer/pkg/telegram/model"
)

type DefaultGateway struct {
	client          externalClient
	commandHandlers map[string]Handler
	messageHandlers []Handler
	middlewares     []Middleware
	log             logger.Logger
}

func NewGateway(c externalClient, log logger.Logger) *DefaultGateway {
	return &DefaultGateway{client: c, log: log}
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
				g.logHandleUpdateError(&request, err)

				_, _ = senderImpl.Send(
					model.
						NewResponse().
						SetChatID(request.Chat.ID).
						SetText("Something went wrong :("),
				)
			}
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

func (g *DefaultGateway) determineCommand(request *model.Request) (Handler, bool) {
	command, ok := g.commandHandlers[request.Command]
	return command, ok
}

func (g *DefaultGateway) logHandleUpdateError(request *model.Request, err error) {
	g.log.Error(
		"telegram handle update error",
		err,
		logger.Field{Key: "user_tg_id", Value: request.User.ID},
		logger.Field{Key: "user_id", Value: request.User.OriginalID},
		logger.Field{Key: "chat_id", Value: request.Chat.ID},
		logger.Field{Key: "tg_update_id", Value: request.UpdateID},
		logger.Field{Key: "tg_command", Value: request.Command},
		logger.Field{Key: "tg_message_id", Value: request.MessageID},
		logger.Field{Key: "tg_message_text", Value: request.Message.Text},
	)
}
