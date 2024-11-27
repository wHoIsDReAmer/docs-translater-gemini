package gemini

// Requests

type GeminiRequest struct {
	Contents []Content `json:"contents"`
}

type Content struct {
	Parts []GeminiPart `json:"parts"`
}

// Responses

type GeminiResponse struct {
	Candidates    []GeminiCandidate `json:"candidates"`
	UsageMetadata UsageMetadata     `json:"usageMetadata"`
	ModelVersion  string            `json:"modelVersion"`
}

type UsageMetadata struct {
	PromptTokenCount     int `json:"promptTokenCount"`
	CandidatesTokenCount int `json:"candidatesTokenCount"`
	TotalTokenCount      int `json:"totalTokenCount"`
}

type GeminiCandidate struct {
	Content struct {
		Parts []GeminiPart `json:"parts"`
		Role  string       `json:"role"`
	} `json:"content"`
	FinishReason     string           `json:"finishReason"`
	CitationMetadata CitationMetadata `json:"citationMetadata"`
	AvgLogProbs      float64          `json:"avgLogprobs"`
}

type GeminiPart struct {
	Text string `json:"text"`
}

type CitationMetadata struct {
	CitationSource []CitiationSource `json:"citationSource"`
}

type CitiationSource struct {
	StartIndex int    `json:"startIndex"`
	EndIndex   int    `json:"endIndex"`
	URI        string `json:"uri"`
}
