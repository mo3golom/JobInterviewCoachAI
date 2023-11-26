package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"job-interviewer/pkg/telegram/model"
)

type sender struct {
	client externalClient
}

func (s *sender) Send(response model.Response) (int64, error) {
	result, err := s.client.Send(response.ToMessageConfig())
	if err != nil {
		return 0, err
	}

	return int64(result.MessageID), nil
}

func (s *sender) Update(messageID int64, response model.Response) error {
	message := response.ToEditMessageTextConfig()
	message.MessageID = int(messageID)

	_, err := s.client.Send(message)
	return err
}

func (s *sender) SendCallback(callbackID model.CallbackID, message ...string) error {
	msg := ""
	if len(message) > 0 {
		msg = message[0]
	}

	_, err := s.client.Request(
		tgbotapi.NewCallback(string(callbackID), msg),
	)
	return err
}
