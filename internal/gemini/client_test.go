package gemini

import "testing"

func TestGenerateText(t *testing.T) {
	client := NewGeminiClient(&GeminiConfig{
		ApiKey: "input_your_gemini_token_here",
	})

	resp, err := client.GenerateText("Explain how the internet works")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(resp)
}
