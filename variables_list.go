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
)
