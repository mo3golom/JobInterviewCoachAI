package ru

import (
	languageService "job-interviewer/internal/telegram/language"
	"job-interviewer/pkg/language"
)

type Dict struct {
}

func (d Dict) GetTexts() map[language.TextKey]string {
	return map[language.TextKey]string{
		languageService.Start: `Привет! Я ваш идеальный помощник в подготовке к собеседованию на разных языках!

Просто нажмите "🆕 Начать новое интервью", выберите специализацию, на которой хотите потренироваться. А дальше я сгенерирую список соответствующих вопросов, используя модель GPT от openAI. 

Как только вы ответите на вопрос, я дам вам конструктивную обратную связь, и буду готов отправить следующий вопрос :)
        `,
		languageService.StartInterview:        "🆕 Начать новое интервью",
		languageService.StartInterviewShort:   "🆕 Новое",
		languageService.ContinueInterview:     "➡️ Продолжить",
		languageService.ActiveInterviewExists: "Хм... у вас уже есть активное интервью! Хотите продолжить или начать новое?",
		languageService.ChoosePosition:        "Для начала выбери позицию, для которой хочешь пройти интервью:",
		languageService.ChooseLevel:           "А теперь выбери свой уровень:",
		languageService.StartInterviewSummary: `
Начинаем интервью!
Позиция: %s
Уровень: %s
        `,
		languageService.LoadQuestions:   "🔄 Нейросеть подготавливает вопросы, одну минутку...",
		languageService.FinishInterview: "️⏏️️ Завершить",
	}
}
