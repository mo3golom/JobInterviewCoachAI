package go_openai

const (
	OpenAI ServiceType = "openai"
	VseGPT ServiceType = "vsegpt"
)

type (
	ServiceType string

	Config struct {
		ServiceType ServiceType
		OpenAI      *OpenAIConfig
		VseGPT      *VseGPTConfig
	}

	VseGPTConfig struct {
		BaseUrl   string
		AuthToken string
	}

	OpenAIConfig struct {
		AuthToken string
	}
)
