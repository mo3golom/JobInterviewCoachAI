package service

import (
	"context"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	interviewerContracts "job-interviewer/internal/interviewer/contracts"
	"job-interviewer/internal/telegram/handlers"
	"job-interviewer/pkg/language"
	"job-interviewer/pkg/telegram"
	"job-interviewer/pkg/telegram/model"
	"job-interviewer/pkg/telegram/service/keyboard"
)

type DefaultService struct {
	finishInterviewUC interviewerContracts.FinishInterviewUseCase
	getNextQuestionUC interviewerContracts.GetNextQuestionUseCase
	languageStorage   language.Storage
}

func NewService(
	finishInterviewUC interviewerContracts.FinishInterviewUseCase,
	getNextQuestionUC interviewerContracts.GetNextQuestionUseCase,
) *DefaultService {
	return &DefaultService{
		finishInterviewUC: finishInterviewUC,
		getNextQuestionUC: getNextQuestionUC,
		languageStorage:   configLanguage(),
	}
}

func (s *DefaultService) FinishInterview(ctx context.Context, request *model.Request, sender telegram.Sender) error {
	userLang := request.User.Lang
	summary, err := s.finishInterviewUC.FinishInterview(ctx, request.User.OriginalID)
	if err != nil {
		return err
	}

	outMessage := s.languageStorage.GetText(userLang, textKeyFinishInterview)
	if summary != "" {
		outMessage = fmt.Sprintf(
			`%s
                    %s`,
			outMessage,
			summary,
		)
	}

	_, err = sender.Send(
		model.NewResponse().
			SetText(
				fmt.Sprintf(
					"%s %s",
					handlers.RobotPrefix,
					outMessage,
				),
			),
	)
	return err
}

func (s *DefaultService) GetNextQuestion(ctx context.Context, request *model.Request, sender telegram.Sender) error {
	userLang := request.User.Lang
	response := model.NewResponse()

	question, err := s.getNextQuestionUC.GetNextQuestion(ctx, request.User.OriginalID)
	if errors.Is(err, interviewerContracts.ErrEmptyActiveInterview) {
		_, err = sender.Send(response.SetText(
			fmt.Sprintf(
				"%s %s",
				handlers.RobotPrefix,
				s.languageStorage.GetText(userLang, textKeyNotFoundActiveInterview),
			),
		))
		return err
	}
	if err != nil {
		return err
	}

	inlineKeyboard, err := keyboard.BuildInlineKeyboardGrid(
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

func (s *DefaultService) GetUserMainKeyboard(lang language.Language) *tgbotapi.ReplyKeyboardMarkup {
	return keyboard.BuildKeyboardGrid(
		keyboard.BuildKeyboardIn{
			Buttons: []keyboard.Button{
				{
					Value: s.languageStorage.GetText(lang, textKeyStartInterview),
				},
			},
		},
	)
}
