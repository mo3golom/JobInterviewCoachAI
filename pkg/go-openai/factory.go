package go_openai

import (
	"github.com/sashabaranov/go-openai"
)

func NewClient(config Config) *openai.Client {
	switch config.ServiceType {
	case VseGPT:
		if config.VseGPT == nil {
			panic("config VseGPT is absent")
		}

		clientConfig := openai.DefaultConfig(config.VseGPT.AuthToken)
		clientConfig.BaseURL = config.VseGPT.BaseUrl
		return openai.NewClientWithConfig(clientConfig)
	default:
		if config.OpenAI == nil {
			panic("config OpenAI is absent")
		}

		return openai.NewClient(
			config.OpenAI.AuthToken,
		)
	}
}
