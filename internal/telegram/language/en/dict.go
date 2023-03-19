package en

import (
	languageService "job-interviewer/internal/telegram/language"
	"job-interviewer/pkg/language"
)

type Dict struct {
}

func (d Dict) GetTexts() map[language.TextKey]string {
	return map[language.TextKey]string{
		languageService.ProcessingAnswer:        "Processing your answer...",
		languageService.NotFoundActiveInterview: "I can`t find an active interview T-T",
		languageService.FinishInterviewSummary:  "Interview’s over! Well done!",
	}
}
