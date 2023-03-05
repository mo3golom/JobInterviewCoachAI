package service

import (
	"context"
	"errors"
	"fmt"
	interviewerContracts "job-interviewer/internal/interviewer/contracts"
	"job-interviewer/internal/telegram/handlers"
	"job-interviewer/internal/telegram/storage"
	"job-interviewer/pkg/telegram"
	"job-interviewer/pkg/telegram/model"
	"job-interviewer/pkg/telegram/service/keyboard"
)

type DefaultService struct {
	finishInterviewUC interviewerContracts.FinishInterviewUseCase
	getNextQuestionUC interviewerContracts.GetNextQuestionUseCase
	keyboardService   keyboard.Service
	storage           storage.Storage
}

func NewService(
	finishInterviewUC interviewerContracts.FinishInterviewUseCase,
	getNextQuestionUC interviewerContracts.GetNextQuestionUseCase,
	keyboardService keyboard.Service,
	storage storage.Storage,
) *DefaultService {
	return &DefaultService{
		finishInterviewUC: finishInterviewUC,
		getNextQuestionUC: getNextQuestionUC,
		keyboardService:   keyboardService,
		storage:           storage,
	}
}

func (s *DefaultService) FinishInterview(ctx context.Context, request *model.Request, sender telegram.Sender) error {
	err := s.finishInterviewUC.FinishInterview(ctx, request.User.OriginalID)
	if err != nil {
		return err
	}

	_, err = sender.Send(
		model.NewResponse(request.Chat.ID).
			SetText(handlers.FinishText),
	)
	if err != nil {
		return err
	}

	return s.HideInlineKeyboardForBotLastMessage(ctx, request, sender)
}

func (s *DefaultService) GetNextQuestion(ctx context.Context, request *model.Request, sender telegram.Sender) error {
	response := model.NewResponse(request.Chat.ID)

	err := s.HideInlineKeyboardForBotLastMessage(ctx, request, sender)
	if err != nil {
		return err
	}

	question, err := s.getNextQuestionUC.GetNextQuestion(ctx, request.User.OriginalID)
	if errors.Is(err, interviewerContracts.ErrNextQuestionEmpty) {
		return s.FinishInterview(ctx, request, sender)
	}
	if errors.Is(err, interviewerContracts.ErrEmptyActiveInterview) {
		_, err = sender.Send(response.SetText(handlers.NoActiveInterviewText))
		return err
	}
	if err != nil {
		return err
	}

	inlineKeyboard, err := s.keyboardService.BuildInlineKeyboardGrid(
		keyboard.BuildInlineKeyboardIn{
			Buttons: getNextQuestionButtons,
		},
	)
	if err != nil {
		return err
	}
	messageID, err := sender.Send(
		response.
			SetText(fmt.Sprintf(handlers.RobotPrefixText, question.Text)).
			SetInlineKeyboardMarkup(inlineKeyboard),
	)
	if err != nil {
		return err
	}

	return s.SaveBotLastMessageID(ctx, request.Chat.ID, messageID)
}

func (s *DefaultService) SaveBotLastMessageID(ctx context.Context, chatID int64, lastBotMessageID int64) error {
	return s.storage.UpsertTelegramBotDetails(
		ctx,
		storage.UpsertTelegramBotDetailsIn{
			ChatID:           chatID,
			LastBotMessageID: lastBotMessageID,
		},
	)
}

func (s *DefaultService) HideInlineKeyboardForBotLastMessage(ctx context.Context, request *model.Request, sender telegram.Sender) error {
	botLastMessageID, err := s.storage.GetBotLastMessageID(ctx, request.Chat.ID)
	if errors.Is(err, storage.ErrEmptyTelegramBotDetailsResult) {
		return nil
	}
	if err != nil {
		return err
	}
	if botLastMessageID == 0 {
		return nil
	}

	messageText := ""
	if request.Message != nil {
		messageText = request.Message.Text
	}

	return sender.Update(
		botLastMessageID,
		model.NewResponse(request.Chat.ID).
			SetText(messageText).
			SetInlineKeyboardMarkup(nil),
	)
}
