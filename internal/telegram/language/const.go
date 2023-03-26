package language

import "job-interviewer/pkg/language"

const (
	English language.Language = "en"
	Russian language.Language = "ru"
)

const (
	Start                   = "START"
	StartInterview          = "START_INTERVIEW"
	StartInterviewShort     = "START_INTERVIEW_SHORT"
	StartInterviewSummary   = "START_INTERVIEW_SUMMARY"
	FinishInterview         = "FINISH_INTERVIEW"
	FinishInterviewSummary  = "FINISH_INTERVIEW_SUMMARY"
	ContinueInterview       = "CONTINUE_INTERVIEW"
	ActiveInterviewExists   = "ACTIVE_INTERVIEW_EXISTS"
	ChoosePosition          = "CHOOSE_POSITION"
	ChooseLevel             = "CHOOSE_LEVEL"
	LoadQuestions           = "LOAD_QUESTIONS"
	ProcessingAnswer        = "PROCESSING_ANSWER"
	NotFoundActiveInterview = "NOT_FOUND_ACTIVE_INTERVIEW"
	ChooseLanguage          = "CHOOSE_LANGUAGE"
	ChooseLanguageSettings  = "CHOOSE_LANGUAGE_SETTINGS"
	ChooseLanguageSuccess   = "CHOOSE_LANGUAGE_SUCCESS"
)
