package ai

type OllamaGenerateRequest struct {
	Model       string  `json:"model"`
	Prompt      string  `json:"prompt"`
	Stream      bool    `json:"stream"`
	Temperature float64 `json:"temperature"`
}

type OllamaGenerateResponse struct {
	Response string `json:"response"`
}