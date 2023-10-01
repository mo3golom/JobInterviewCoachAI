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
	userLang := language.Russian

	messageID, err := sender.Send(
		model.NewResponse().
			SetText(
				fmt.Sprintf(
					"%s %s",
					handlers.RobotPrefix,
					"loading...",
				),
			))
	if err != nil {
		return err
	}

	summary, err := s.finishInterviewUC.FinishInterview(ctx, request.User.OriginalID)
	if err != nil {
		return err
	}

	outMessage := s.languageStorage.GetText(userLang, textKeyFinishInterview)
	if summary != "" {
		outMessage = fmt.Sprintf(
			`%s %s`,
			outMessage,
			summary,
		)
	}

	return sender.Update(
		messageID,
		model.NewResponse().
			SetText(
				fmt.Sprintf(
					"%s %s",
					handlers.RobotPrefix,
					outMessage,
				),
			),
	)
}

func (s *DefaultService) GetNextQuestion(ctx context.Context, request *model.Request, sender telegram.Sender, updateMessageID ...int64) error {
	userLang := request.User.Lang
	var targetUpdateMessageID int64
	if len(updateMessageID) > 0 {
		targetUpdateMessageID = updateMessageID[0]
	} else {
		var err error
		targetUpdateMessageID, err = sender.Send(
			model.NewResponse().
				SetText(
					fmt.Sprintf(
						"%s %s",
						handlers.RobotPrefix,
						"loading...",
					),
				))
		if err != nil {
			return err
		}
	}

	question, err := s.getNextQuestionUC.GetNextQuestion(ctx, request.User.OriginalID)
	if errors.Is(err, interviewerContracts.ErrEmptyActiveInterview) {
		return sender.Update(
			targetUpdateMessageID,
			model.NewResponse().
				SetText(
					fmt.Sprintf(
						"%s %s",
						handlers.RobotPrefix,
						s.languageStorage.GetText(userLang, textKeyNotFoundActiveInterview),
					),
				))
	}
	if errors.Is(err, interviewerContracts.ErrQuestionsInFreePlanHaveExpired) {
		err = sender.Update(
			targetUpdateMessageID,
			model.NewResponse().
				SetText(
					fmt.Sprintf(
						"%s %s",
						handlers.RobotPrefix,
						s.languageStorage.GetText(language.Russian, textKeyFreeQuestionsIsEnd),
					),
				))

		if err != nil {
			return err
		}

		return s.FinishInterview(ctx, request, sender)
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

	return sender.Update(
		targetUpdateMessageID,
		model.NewResponse().
			SetText(fmt.Sprintf("%s %s", handlers.RobotPrefix, question.Text)).
			SetInlineKeyboardMarkup(inlineKeyboard),
	)
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

func (s *DefaultService) ShowSubscribeMessage(sender telegram.Sender) error {
	userLang := language.Russian

	inlineKeyboard, err := keyboard.BuildInlineKeyboardGrid(
		keyboard.BuildInlineKeyboardIn{
			Buttons: subscribeButtons,
		},
	)
	if err != nil {
		return err
	}

	_, err = sender.Send(
		model.NewResponse().
			SetText(s.languageStorage.GetText(userLang, textKeySubscribe)).
			SetInlineKeyboardMarkup(inlineKeyboard),
	)
	return err
}
