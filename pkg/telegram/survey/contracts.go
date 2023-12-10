package survey

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type (
	Survey[D any] interface {
		FindUnansweredQuestionAsKeyboard(command string) (*QuestionWithKeyboard, error)
		Init(in []string, out D) Survey[D]
		AddQuestion(question Question[D], additional ...Question[D]) Survey[D]
		AddQuestionWhen(question Question[D], clause func() bool) Survey[D]
	}

	Question[D any] interface {
		toInlineKeyboard(command string, previousAnswers ...string) (*tgbotapi.InlineKeyboardMarkup, error)
		text() string
		isAnswered() bool
		answerID() int
		setAnswer(answerID int, out D) Question[D]
	}

	QuestionWithKeyboard struct {
		Text     string
		Keyboard *tgbotapi.InlineKeyboardMarkup
	}
)
