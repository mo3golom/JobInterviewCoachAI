package startinterview

import (
	"job-interviewer/pkg/language"
	"strconv"
)

const (
	textKeyStartInterview     language.TextKey = "textKeyStartInterview"
	textKeyClarifyJobPosition language.TextKey = "textKeyClarifyJobPosition"
	textKeyDeveloper          language.TextKey = "developer"
	textKeyProjectManager     language.TextKey = "project_manager"
	textKeyProductManager     language.TextKey = "product_manager"
	textKeyProductDesigner    language.TextKey = "product_designer"
	textKeyGolang             language.TextKey = "golang"
	textKeyPHP                language.TextKey = "php"
	textKeyPython             language.TextKey = "python"
	textKeyRust               language.TextKey = "rust"
	textKeyJavascript         language.TextKey = "javascript"
	textKeySwift              language.TextKey = "swift"
	textKeyJava               language.TextKey = "java"
	textKeyCPlusPlus          language.TextKey = "c_plus_plus"
	textKeyCSharp             language.TextKey = "c_sharp"
	textKeyYes                language.TextKey = "yes"
	textKeyNo                 language.TextKey = "no"
)

func configLanguage() language.Storage {
	return language.NewLangStorage(
		map[language.Language]language.WordStorage{
			language.Russian: language.NewWordStorage(
				map[language.TextKey]string{
					language.TextKey(strconv.FormatInt(QuestionContinueActiveInterview, 10)): "У вас есть уже активное интервью %s! Хотите продолжить?",
					language.TextKey(strconv.FormatInt(QuestionJobPosition, 10)):             "Выбери позицию, для которой хочешь пройти интервью:",
					textKeyClarifyJobPosition: "Ты выбрал: \"%s\", осталось уточнить направление:",
					textKeyStartInterview:     "🚀 Начать новое интервью",
					textKeyDeveloper:          "Разработчик",
					textKeyProjectManager:     "Project менеджер",
					textKeyProductManager:     "Product менеджер",
					textKeyProductDesigner:    "Product дизайнер",
					textKeyGolang:             "Golang",
					textKeyPHP:                "PHP",
					textKeyPython:             "Python",
					textKeyRust:               "Rust",
					textKeyJavascript:         "Javascript",
					textKeySwift:              "Swift",
					textKeyJava:               "Java",
					textKeyCPlusPlus:          "C++",
					textKeyCSharp:             "C#",
					textKeyYes:                "Да",
					textKeyNo:                 "Нет",
				},
			),
		},
	)
}
