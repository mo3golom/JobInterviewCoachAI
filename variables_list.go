package job_interviewer

import (
	"job-interviewer/pkg/variables"
)

var (
	GPTApiKey                = variables.Environment[string]("GPT_API_KEY", "")
	TGBotToken               = variables.Environment[string]("TG_BOT_TOKEN", "")
	MonthlySubscriptionPrice = variables.Environment[string]("MONTHLY_SUBSCRIPTION_PRICE", "99")
	YMShopID                 = variables.Environment[string]("YM_SHOP_ID", "")
	YMSecretKey              = variables.Environment[string]("YM_SECRET_KEY", "")
	PaidModelEnable          = variables.Environment[bool]("PAID_MODEL_ENABLE", true)
	FreeQuestionsCount       = variables.Environment[string]("FREE_QUESTIONS_COUNT", "5")
	FreeInterviewsCount      = variables.Environment[string]("FREE_INTERVIEWS_COUNT", "10")
)
