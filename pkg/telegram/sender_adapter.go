package telegram

import (
	"job-interviewer/pkg/telegram/model"
)

type senderAdapter struct {
	*sender
	chatID int64
}

func (s senderAdapter) Send(response model.Response) (int64, error) {
	response = response.SetChatID(s.chatID)
	return s.sender.Send(response)
}

func (s senderAdapter) Update(messageID int64, response model.Response) error {
	response = response.SetChatID(s.chatID)
	return s.sender.Update(messageID, response)
}
