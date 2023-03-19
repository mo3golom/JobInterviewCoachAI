package service

import (
	"context"
	"errors"
	"fmt"
	interviewerContracts "job-interviewer/internal/interviewer/contracts"
	"job-interviewer/internal/telegram/handlers"
	languageService "job-interviewer/internal/telegram/language"
	"job-interviewer/internal/telegram/storage"
	"job-interviewer/pkg/language"
	"job-interviewer/pkg/telegram"
	"job-interviewer/pkg/telegram/model"
	"job-interviewer/pkg/telegram/service/keyboard"
)

type DefaultService struct {
	finishInterviewUC interviewerContracts.FinishInterviewUseCase
	getNextQuestionUC interviewerContracts.GetNextQuestionUseCase
	keyboardService   keyboard.Service
	storage           storage.Storage
	languageService   languageService.Service
}

func NewService(
	finishInterviewUC interviewerContracts.FinishInterviewUseCase,
	getNextQuestionUC interviewerContracts.GetNextQuestionUseCase,
	keyboardService keyboard.Service,
	storage storage.Storage,
	l languageService.Service,
) *DefaultService {
	return &DefaultService{
		finishInterviewUC: finishInterviewUC,
		getNextQuestionUC: getNextQuestionUC,
		keyboardService:   keyboardService,
		storage:           storage,
		languageService:   l,
	}
}

func (s *DefaultService) Start(request *model.Request, sender telegram.Sender) error {
	userLang := language.Language(request.User.Lang)

	_, err := sender.Send(
		model.NewResponse(request.Chat.ID).
			SetText(s.languageService.GetText(userLang, languageService.Start)).
			SetKeyboardMarkup(
				s.keyboardService.BuildKeyboardGrid(
					keyboard.BuildKeyboardIn{
						Buttons: []keyboard.Button{
							{
								Value: s.languageService.GetText(userLang, languageService.StartInterview),
							},
						},
					},
				),
			),
	)

	return err
}

func (s *DefaultService) FinishInterview(ctx context.Context, request *model.Request, sender telegram.Sender) error {
	err := s.finishInterviewUC.FinishInterview(ctx, request.User.OriginalID)
	if err != nil {
		return err
	}

	_, err = sender.Send(
		model.NewResponse(request.Chat.ID).
			SetText(
				fmt.Sprintf(
					"%s %s",
					handlers.RobotPrefix,
					s.languageService.GetText(languageService.English, languageService.FinishInterviewSummary),
				),
			),
	)
	return err
}

func (s *DefaultService) GetNextQuestion(ctx context.Context, request *model.Request, sender telegram.Sender) error {
	response := model.NewResponse(request.Chat.ID)

	question, err := s.getNextQuestionUC.GetNextQuestion(ctx, request.User.OriginalID)
	if errors.Is(err, interviewerContracts.ErrNextQuestionEmpty) {
		return s.FinishInterview(ctx, request, sender)
	}
	if errors.Is(err, interviewerContracts.ErrEmptyActiveInterview) {
		_, err = sender.Send(response.SetText(
			fmt.Sprintf(
				"%s %s",
				handlers.RobotPrefix,
				s.languageService.GetText(languageService.English, languageService.NotFoundActiveInterview),
			),
		))
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
	_, err = sender.Send(
		response.
			SetText(fmt.Sprintf("%s %s", handlers.RobotPrefix, question.Text)).
			SetInlineKeyboardMarkup(inlineKeyboard),
	)

	return err
}

// SaveBotLastMessageID Пока не используется, потому что могут быть интересные баги
func (s *DefaultService) SaveBotLastMessageID(ctx context.Context, chatID int64, lastBotMessageID int64) error {
	return s.storage.UpsertTelegramBotDetails(
		ctx,
		storage.UpsertTelegramBotDetailsIn{
			ChatID:           chatID,
			LastBotMessageID: lastBotMessageID,
		},
	)
}

// HideInlineKeyboardForBotLastMessage Пока не используется, потому что могут быть интересные баги
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
