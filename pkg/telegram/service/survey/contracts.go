package survey

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type (
	Survey interface {
		FindUnansweredQuestionAsKeyboard(command string) (*QuestionWithKeyboard, error)
		SetAnswers(requestData []string) Survey
		AddQuestion(question Question, additional ...Question) Survey
	}

	Question interface {
		toInlineKeyboard(command string, previousAnswers ...string) (*tgbotapi.InlineKeyboardMarkup, error)
		text() string
		isAnswered() bool
		answerID() int
		setAnswer(answerID int) Question
	}

	QuestionWithKeyboard struct {
		Text     string
		Keyboard *tgbotapi.InlineKeyboardMarkup
	}
)
