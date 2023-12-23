package survey

import (
	"strconv"
)

type (
	DefaultSurvey[D any] struct {
		questions []Question[D]
	}
)

func New[D any]() Survey[D] {
	return DefaultSurvey[D]{
		questions: make([]Question[D], 0, 1),
	}
}

func (s DefaultSurvey[D]) FindUnansweredQuestionAsKeyboard(command string) (*QuestionWithKeyboard, error) {
	var q Question[D]
	previousAnswers := make([]string, 0, len(s.questions))
	for _, question := range s.questions {
		q = question
		if q.isAnswered() {
			previousAnswers = append(previousAnswers, strconv.Itoa(question.answerID()))

			continue
		}

		keyboard, err := q.toInlineKeyboard(command, previousAnswers...)
		return &QuestionWithKeyboard{
			Text:     q.text(),
			Keyboard: keyboard,
		}, err
	}

	return nil, nil
}

func (s DefaultSurvey[D]) Init(requestData []string, out D) Survey[D] {
	for index, data := range requestData {
		answerID, err := strconv.Atoi(data)
		if err != nil {
			continue
		}

		if index >= len(s.questions) {
			break
		}

		s.questions[index] = s.questions[index].setAnswer(answerID, out)
	}

	return s
}

func (s DefaultSurvey[D]) AddQuestion(question Question[D], additional ...Question[D]) Survey[D] {
	s.questions = append(s.questions, question)

	if len(additional) > 0 {
		s.questions = append(s.questions, additional...)
	}

	return s
}

func (s DefaultSurvey[D]) AddQuestionWhen(question Question[D], clause func() bool) Survey[D] {
	if !clause() {
		return s
	}

	return s.AddQuestion(question)
}
