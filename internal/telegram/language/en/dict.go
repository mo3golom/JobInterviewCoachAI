package en

import (
	languageService "job-interviewer/internal/telegram/language"
	"job-interviewer/pkg/language"
)

type Dict struct {
}

func (d Dict) GetTexts() map[language.TextKey]string {
	return map[language.TextKey]string{
		languageService.Start: `Hello! I am your perfect assistant in preparing for interviews in different languages!

Simply click "🆕 Start a new interview", choose the specialization you want to practice, and then I will generate a list of corresponding questions using the openAI GPT model.

As soon as you answer a question, I will provide you with constructive feedback and be ready to send the next question :)
        `,
		languageService.StartInterview:        "🆕 Start a new interview",
		languageService.StartInterviewShort:   "🆕 New",
		languageService.ContinueInterview:     "➡️ Continue",
		languageService.ActiveInterviewExists: "Hmm... you already have an active interview! Do you want to continue or start a new one?",
		languageService.ChoosePosition:        "To begin, choose the position for which you want to take the interview:",
		languageService.ChooseLevel:           "And now, choose your level:",
		languageService.StartInterviewSummary: `
Let's start!
Position: %s
        `,
		languageService.ProcessingAnswer:        "Processing your answer...",
		languageService.NotFoundActiveInterview: "I can`t find an active interview T-T",
		languageService.FinishInterviewSummary:  "Interview’s over! Well done!",
		languageService.ChooseLanguage:          "Choose bot ui language",
		languageService.ChooseLanguageSettings:  "⚙️ Change language",
		languageService.ChooseLanguageSuccess:   "Language changed successfully!",
	}
}
