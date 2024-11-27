package gemini

import "time"

type GeminiConfig struct {
	ApiKey  string
	Proxy   string
	Timeout time.Duration
}

func NewGeminiConfig(apiKey string) *GeminiConfig {
	return &GeminiConfig{ApiKey: apiKey}
}
