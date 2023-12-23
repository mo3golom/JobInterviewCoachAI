package survey

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	keyboard2 "job-interviewer/pkg/telegram/keyboard"
	"sort"
	"strconv"
	"unicode/utf8"
)

const (
	maxKeyboardButtonContentLen = 10
)

type (
	PossibleAnswer[T comparable] struct {
		content string
		value   T
	}

	DefaultQuestion[T comparable, D any] struct {
		content         string
		possibleAnswers []PossibleAnswer[T]
		setAnswerFunc   func(answer T, out D)
		answerIDValue   int
		isAnsweredValue bool
	}
)

func NewPossibleAnswer(content string) PossibleAnswer[string] {
	answer := PossibleAnswer[string]{
		content: content,
		value:   content,
	}

	return answer
}

func NewComplexPossibleAnswer[T comparable](content string, value ...T) PossibleAnswer[T] {
	answer := PossibleAnswer[T]{
		content: content,
	}

	if len(value) > 0 {
		answer.value = value[0]
	}

	return answer
}

func NewQuestion[T comparable, D any](
	question string,
	setAnswerFunc func(answer T, out D),
	possibleAnswers ...PossibleAnswer[T],
) Question[D] {
	return DefaultQuestion[T, D]{
		content:         question,
		possibleAnswers: possibleAnswers,
		setAnswerFunc:   setAnswerFunc,
	}
}

func (q DefaultQuestion[T, D]) toInlineKeyboard(command string, previousAnswers ...string) (*tgbotapi.InlineKeyboardMarkup, error) {
	keyboardListView := false
	buttons := make([]keyboard2.InlineButton, 0, len(q.possibleAnswers))
	for index, value := range q.possibleAnswers {
		data := make([]string, 0, 1)
		data = append(data, previousAnswers...)
		data = append(data, strconv.Itoa(index))
		buttons = append(
			buttons,
			keyboard2.InlineButton{
				Value: value.content,
				Data:  data,
				Type:  keyboard2.ButtonData,
			},
		)

		if utf8.RuneCountInString(value.content) > maxKeyboardButtonContentLen {
			keyboardListView = true
		}
	}

	sort.Slice(buttons, func(i, j int) bool {
		return buttons[i].Value < buttons[j].Value
	})

	if keyboardListView {
		return keyboard2.BuildInlineKeyboardList(
			keyboard2.BuildInlineKeyboardIn{
				Command: &command,
				Buttons: buttons,
			},
		)
	}
	return keyboard2.BuildInlineKeyboardGrid(
		keyboard2.BuildInlineKeyboardIn{
			Command: &command,
			Buttons: buttons,
		},
	)
}

func (q DefaultQuestion[T, D]) text() string {
	return q.content
}

func (q DefaultQuestion[T, D]) isAnswered() bool {
	return q.isAnsweredValue
}

func (q DefaultQuestion[T, D]) answerID() int {
	return q.answerIDValue
}

func (q DefaultQuestion[T, D]) setAnswer(answerID int, out D) Question[D] {
	q.isAnsweredValue = true
	q.answerIDValue = answerID
	q.setAnswerFunc(q.possibleAnswers[answerID].value, out)

	return q
}

func (p PossibleAnswer[T]) GetContent() string {
	return p.content
}
