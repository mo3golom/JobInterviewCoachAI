package telegram

import (
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
