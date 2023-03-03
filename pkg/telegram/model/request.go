package model

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"strings"
)

type (
	Request struct {
		UpdateID  int64
		MessageID int64
		Command   string
		Data      []string
		Chat      *Chat
		User      *User
		Message   *Message
	}

	Chat struct {
		ID int64
	}

	User struct {
		ID         int64
		OriginalID uuid.UUID
	}

	Message struct {
		Text string
	}
)

func NewRequest(in tgbotapi.Update) Request {
	chat := in.FromChat()
	user := in.SentFrom()
	request := Request{
		UpdateID: int64(in.UpdateID),
		Chat: &Chat{
			ID: chat.ID,
		},
		User: &User{
			ID: user.ID,
		},
	}

	if in.Message != nil && len(in.Message.Text) <= 150 {
		request.MessageID = int64(in.Message.MessageID)
		request.Message = &Message{
			Text: in.Message.Text,
		}

		request.Command = in.Message.Text
		request.Data = []string{}
	}

	if in.CallbackQuery == nil {
		return request
	}
	request.MessageID = int64(in.CallbackQuery.Message.MessageID)
	command, data := StringToCommandWithData(in.CallbackQuery.Data)
	request.Command = command
	request.Data = data

	return request
}

func (r *Request) toCallbackString() string {
	data := strings.Join(r.Data, ":")
	return fmt.Sprintf("%s#%s", r.Command, data)
}
