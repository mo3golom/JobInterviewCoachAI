package service

import (
	"context"
	"errors"
	"fmt"
	interviewerContracts "job-interviewer/internal/interviewer/contracts"
	"job-interviewer/internal/telegram/handlers"
	"job-interviewer/pkg/telegram"
	"job-interviewer/pkg/telegram/model"
	"job-interviewer/pkg/telegram/service/keyboard"
)

type DefaultService struct {
	finishInterviewUC interviewerContracts.FinishInterviewUseCase
	getNextQuestionUC interviewerContracts.GetNextQuestionUseCase
	keyboardService   keyboard.Service
}

func NewService(
	finishInterviewUC interviewerContracts.FinishInterviewUseCase,
	getNextQuestionUC interviewerContracts.GetNextQuestionUseCase,
	keyboardService keyboard.Service,
) *DefaultService {
	return &DefaultService{
		finishInterviewUC: finishInterviewUC,
		getNextQuestionUC: getNextQuestionUC,
		keyboardService:   keyboardService,
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

	return nil
}

func (s *DefaultService) GetNextQuestion(ctx context.Context, request *model.Request, sender telegram.Sender) error {
	response := model.NewResponse(request.Chat.ID)

	question, err := s.getNextQuestionUC.GetNextQuestion(ctx, request.User.OriginalID)
	if errors.Is(err, interviewerContracts.ErrNextQuestionEmpty) {
		return s.FinishInterview(ctx, request, sender)
	}
	if errors.Is(err, interviewerContracts.ErrEmptyActiveInterview) {
		_, err = sender.Send(response.SetText(handlers.NoActiveInterviewText))
		if err != nil {
			return err
		}
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
			SetText(fmt.Sprintf(handlers.RobotPrefixText, question.Text)).
			SetInlineKeyboardMarkup(inlineKeyboard),
	)
	if err != nil {
		return err
	}

	return nil
}

// TODO: необзодимо запоминать последнее сообщение бота  и у него скрывать клавиатуру
func (s *DefaultService) HideInlineKeyboard(request *model.Request, sender telegram.Sender) error {
	messageText := ""
	if request.Message != nil {
		messageText = request.Message.Text
	}

	return sender.Update(
		request.MessageID,
		model.NewResponse(request.Chat.ID).
			SetText(messageText).
			SetInlineKeyboardMarkup(nil),
	)
}
