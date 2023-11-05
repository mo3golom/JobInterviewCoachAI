package survey

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"job-interviewer/pkg/telegram/service/keyboard"
	"sort"
	"strconv"
)

type (
	PossibleAnswer[T comparable] struct {
		content string
		value   T
	}

	DefaultQuestion[T comparable] struct {
		content         string
		possibleAnswers []PossibleAnswer[T]
		setAnswerFunc   func(answer T)
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

func NewQuestion[T comparable](
	question string,
	setAnswerFunc func(answer T),
	possibleAnswers ...PossibleAnswer[T],
) Question {
	answers := make([]PossibleAnswer[T], 0, 1+len(possibleAnswers))
	if len(possibleAnswers) > 0 {
		answers = append(answers, possibleAnswers...)
	} else {
		panic("question should have one or more answers!")
	}

	return DefaultQuestion[T]{
		content:         question,
		possibleAnswers: answers,
		setAnswerFunc:   setAnswerFunc,
	}
}

func (q DefaultQuestion[T]) toInlineKeyboard(command string, previousAnswers ...string) (*tgbotapi.InlineKeyboardMarkup, error) {
	buttons := make([]keyboard.InlineButton, 0, len(q.possibleAnswers))
	for index, value := range q.possibleAnswers {
		data := make([]string, 0, 1)
		data = append(data, previousAnswers...)
		data = append(data, strconv.Itoa(index))
		buttons = append(
			buttons,
			keyboard.InlineButton{
				Value: value.content,
				Data:  data,
				Type:  keyboard.ButtonData,
			},
		)
	}

	sort.Slice(buttons, func(i, j int) bool {
		return buttons[i].Value < buttons[j].Value
	})
	return keyboard.BuildInlineKeyboardGrid(
		keyboard.BuildInlineKeyboardIn{
			Command: &command,
			Buttons: buttons,
		},
	)
}

func (q DefaultQuestion[T]) text() string {
	return q.content
}

func (q DefaultQuestion[T]) isAnswered() bool {
	return q.isAnsweredValue
}

func (q DefaultQuestion[T]) answerID() int {
	return q.answerIDValue
}

func (q DefaultQuestion[T]) setAnswer(answerID int) Question {
	q.isAnsweredValue = true
	q.answerIDValue = answerID
	q.setAnswerFunc(q.possibleAnswers[answerID].value)

	return q
}

func (p PossibleAnswer[T]) GetContent() string {
	return p.content
}
