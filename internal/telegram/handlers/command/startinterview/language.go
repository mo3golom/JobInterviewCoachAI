package startinterview

import (
	"job-interviewer/internal"
	"job-interviewer/pkg/language"
)

const (
	textKeyStartInterview                  language.TextKey = "textKeyStartInterview"
	textKeyClarifyJobPosition              language.TextKey = "textKeyClarifyJobPosition"
	textKeyYes                             language.TextKey = "yes"
	textKeyNo                              language.TextKey = "no"
	textKeyQuestionJobPosition             language.TextKey = "questionJobPosition"
	textKeyQuestionContinueActiveInterview language.TextKey = "questionContinueActiveInterview"
	textKeyYourChoice                      language.TextKey = "yourChoice"
)

func configLanguage() language.Storage {
	return language.NewLangStorage(
		map[language.Language]language.WordStorage{
			language.Russian: language.NewWordStorage(
				map[language.TextKey]string{
					textKeyQuestionContinueActiveInterview:         `У вас есть уже активное интервью <b>"%s"</b>! Хотите продолжить?`,
					textKeyQuestionJobPosition:                     "Выбери позицию, для которой хочешь пройти интервью:",
					textKeyClarifyJobPosition:                      `Ты выбрал: <b>"%s"</b>, осталось уточнить направление:`,
					textKeyYourChoice:                              `Ты выбрал: <b>"%s"</b>, тренировка начинается! 🚀`,
					textKeyStartInterview:                          "🚀 Начать новое интервью",
					language.TextKey(Developer):                    "Разработчик",
					language.TextKey(internal.ProjectManager):      "Project менеджер",
					language.TextKey(internal.ProductManager):      "Product менеджер",
					language.TextKey(internal.ProductDesigner):     "Product дизайнер",
					language.TextKey(internal.QAEngineer):          "QA инженер",
					language.TextKey(internal.Behavioral):          "Behavioral интервью",
					language.TextKey(internal.GolangDeveloper):     "Golang",
					language.TextKey(internal.PhpDeveloper):        "PHP",
					language.TextKey(internal.PythonDeveloper):     "Python",
					language.TextKey(internal.RustDeveloper):       "Rust",
					language.TextKey(internal.JavascriptDeveloper): "Javascript",
					language.TextKey(internal.SwiftDeveloper):      "Swift",
					language.TextKey(internal.JavaDeveloper):       "Java",
					language.TextKey(internal.CplusplusDeveloper):  "C++",
					language.TextKey(internal.CsharpDeveloper):     "C#",
					textKeyYes: "Да",
					textKeyNo:  "Нет",
				},
			),
		},
	)
}
