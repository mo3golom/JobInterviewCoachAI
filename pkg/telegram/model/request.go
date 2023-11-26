package model

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"job-interviewer/pkg/language"
)

var (
	lang = map[string]language.Language{
		"ru": language.Russian,
		"en": language.English,
	}
)

type (
	CallbackID string

	Request struct {
		UpdateID   int64
		MessageID  int64
		Command    string
		Data       []string
		CallbackID *CallbackID
		Chat       *Chat
		User       *User
		Message    *Message
		Voice      *Voice
	}

	Chat struct {
		ID int64
	}

	User struct {
		ID         int64
		OriginalID uuid.UUID
		Lang       language.Language
		Username   string
		FirstName  string
		LastName   string
	}

	Message struct {
		Text string
	}

	Voice struct {
		FileID   string
		Duration int
		URL      string
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
			ID:        user.ID,
			Lang:      lang[user.LanguageCode],
			Username:  user.UserName,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		},
	}

	if in.Message != nil {
		request.MessageID = int64(in.Message.MessageID)
		request.Message = &Message{
			Text: in.Message.Text,
		}

		request.Command = in.Message.Text
		request.Data = []string{}

		if in.Message.Voice != nil {
			request.Voice = &Voice{
				FileID:   in.Message.Voice.FileID,
				Duration: in.Message.Voice.Duration,
			}
		}
	}

	if in.CallbackQuery == nil {
		return request
	}
	request.CallbackID = (*CallbackID)(&in.CallbackQuery.ID)
	request.MessageID = int64(in.CallbackQuery.Message.MessageID)
	command, data := StringToCommandWithData(in.CallbackQuery.Data)
	request.Command = command
	request.Data = data
	if in.CallbackQuery.Message != nil {
		request.Message = &Message{
			Text: in.CallbackQuery.Message.Text,
		}
	}

	return request
}
