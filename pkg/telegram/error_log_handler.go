package telegram

import (
	"context"
	"job-interviewer/pkg/logger"
	"job-interviewer/pkg/telegram/model"
)

type ErrorLogHandler struct {
	log logger.Logger
}

func (e *ErrorLogHandler) Handle(_ context.Context, err error, request *model.Request, sender Sender) bool {
	e.log.Error(
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

	_, _ = sender.Send(
		model.
			NewResponse().
			SetText("Something went wrong :("),
	)

	return true
}
