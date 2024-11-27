package gemini

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	GeminiEndpoint = "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent"
)

type GeminiClient struct {
	config *GeminiConfig
}

func NewGeminiClient(config *GeminiConfig) *GeminiClient {
	return &GeminiClient{config: config}
}

func (c *GeminiClient) createHttpClient() (*http.Client, error) {
	client := &http.Client{}

	if c.config.Proxy != "" {
		proxy, err := url.Parse(c.config.Proxy)
		if err != nil {
			return nil, err
		}

		client.Transport = &http.Transport{Proxy: http.ProxyURL(proxy)}
	}

	if c.config.Timeout > 0 {
		client.Timeout = c.config.Timeout
	}
	return client, nil
}

func (c *GeminiClient) GenerateText(prompt string) (string, error) {
	client, err := c.createHttpClient()
	if err != nil {
		return "", err
	}

	request := GeminiRequest{
		Contents: []Content{
			{
				Parts: []GeminiPart{
					{Text: prompt},
				},
			},
		},
	}
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		return "", err
	}
	resp, err := client.Post(fmt.Sprintf("%s?key=%s", GeminiEndpoint, c.config.ApiKey), "application/json", bytes.NewBuffer(jsonRequest))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
