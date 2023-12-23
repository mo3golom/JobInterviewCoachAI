package service

import (
	"context"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	jobInterviewer "job-interviewer"
	interviewerContracts "job-interviewer/internal/interviewer/contracts"
	"job-interviewer/internal/telegram/handlers"
	"job-interviewer/internal/telegram/handlers/command"
	"job-interviewer/pkg/language"
	"job-interviewer/pkg/telegram"
	"job-interviewer/pkg/telegram/keyboard"
	"job-interviewer/pkg/telegram/model"
	"job-interviewer/pkg/variables"
)

type DefaultService struct {
	finishInterviewUC interviewerContracts.FinishInterviewUseCase
	getNextQuestionUC interviewerContracts.GetNextQuestionUseCase
	languageStorage   language.Storage
	variables         variables.Repository
}

func NewService(
	finishInterviewUC interviewerContracts.FinishInterviewUseCase,
	getNextQuestionUC interviewerContracts.GetNextQuestionUseCase,
	variables variables.Repository,
) *DefaultService {
	return &DefaultService{
		finishInterviewUC: finishInterviewUC,
		getNextQuestionUC: getNextQuestionUC,
		variables:         variables,
		languageStorage:   configLanguage(),
	}
}

func (s *DefaultService) FinishInterview(ctx context.Context, request *model.Request, sender telegram.Sender, updateMessageID ...int64) error {
	userLang := language.Russian
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

	summary, err := s.finishInterviewUC.FinishInterview(ctx, request.User.OriginalID)
	if errors.Is(err, interviewerContracts.ErrEmptyActiveInterview) {
		return sender.Update(
			targetUpdateMessageID,
			model.NewResponse().
				SetText(s.languageStorage.GetText(userLang, textKeyFinishNoActiveInterview)),
		)
	}
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
		targetUpdateMessageID,
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
		if !s.variables.GetBool(jobInterviewer.PaidModelEnable) {
			return s.FinishInterview(ctx, request, sender, targetUpdateMessageID)
		}

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
			Buttons: []keyboard.InlineButton{
				{
					Value: s.languageStorage.GetText(userLang, textKeyImDone),
					Data:  []string{command.FinishInterviewCommand},
					Type:  keyboard.ButtonData,
				},
				{
					Value: s.languageStorage.GetText(userLang, textKeySuggestion),
					Data:  []string{command.GetAnswerSuggestionCommand},
					Type:  keyboard.ButtonData,
				},
				{
					Value: s.languageStorage.GetText(userLang, textKeySkipQuestion),
					Data:  []string{command.SkipQuestionCommand},
					Type:  keyboard.ButtonData,
				},
			},
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
	return keyboard.BuildKeyboardCustomGrid(
		keyboard.BuildKeyboardCustomIn{
			Buttons: [][]keyboard.Button{
				{
					{Value: s.languageStorage.GetText(lang, textKeyStartInterview)},
				},
				{
					{Value: s.languageStorage.GetText(lang, textKeySubscription)},
					{Value: s.languageStorage.GetText(lang, textKeyAbout)},
				},
			},
		},
	)
}

func (s *DefaultService) ShowSubscribeMessage(sender telegram.Sender) error {
	if !s.variables.GetBool(jobInterviewer.PaidModelEnable) {
		return nil
	}

	userLang := language.Russian

	inlineKeyboard, err := keyboard.BuildInlineKeyboardGrid(
		keyboard.BuildInlineKeyboardIn{
			Buttons: []keyboard.InlineButton{
				{
					Value: s.languageStorage.GetText(userLang, textKeyBuySubscription),
					Data:  []string{command.PaySubscriptionCommand},
					Type:  keyboard.ButtonData,
				},
			},
		},
	)
	if err != nil {
		return err
	}

	subscriptionPrice := s.variables.GetInt64(jobInterviewer.MonthlySubscriptionPrice)
	_, err = sender.Send(
		model.NewResponse().
			SetText(
				fmt.Sprintf(
					s.languageStorage.GetText(userLang, textKeySubscribe),
					subscriptionPrice,
					s.languageStorage.GetText(userLang, textKeyBuySubscription),
				),
			).
			SetInlineKeyboardMarkup(inlineKeyboard),
	)
	return err
}
