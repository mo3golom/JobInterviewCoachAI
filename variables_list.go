package job_interviewer

import (
	"job-interviewer/pkg/variables"
)

var (
	GPTApiKey  = variables.Environment[string]("GPT_API_KEY", "")
	TGBotToken = variables.Environment[string]("TG_BOT_TOKEN", "")
)
