package survey

import (
	"strconv"
)

type (
	DefaultSurvey struct {
		questions []Question
	}
)

func New() Survey {
	return DefaultSurvey{
		questions: make([]Question, 0, 1),
	}
}

func (s DefaultSurvey) FindUnansweredQuestionAsKeyboard(command string) (*QuestionWithKeyboard, error) {
	var q Question
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

func (s DefaultSurvey) SetAnswers(requestData []string) Survey {
	for index, data := range requestData {
		answerID, err := strconv.Atoi(data)
		if err != nil {
			continue
		}

		if index >= len(s.questions) {
			break
		}

		s.questions[index] = s.questions[index].setAnswer(answerID)
	}

	return s
}

func (s DefaultSurvey) AddQuestion(question Question, additional ...Question) Survey {
	s.questions = append(s.questions, question)

	if len(additional) > 0 {
		s.questions = append(s.questions, additional...)
	}

	return s
}
