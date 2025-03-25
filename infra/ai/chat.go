package ai

import (
	"encoding/json"
	"strings"
)

const (
	CHAT_API_PATH       = "/chat/completions"
	OPENAI_CHAT_POSTFIX = "/v1" + CHAT_API_PATH
)

type ChatCompletionRequest = ChatRequest
type ChatCompletionResponse = ChatResponse

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Stream      bool      `json:"stream"`
	MaxTokens   int32     `json:"max_tokens"`
	Temperature float64   `json:"temperature,omitempty"`
}

type ChatResponse struct {
	Created     int32  `json:"created"`
	Model       string `json:"model"`
	FingerPrint string `json:"system_fingerprint"`
	Choices     []struct {
		// stream=true 时, 使用此字段
		Delta Message `json:"delta"`
		// stream=false 时, 使用此字段
		Message Message `json:"message"`
	} `json:"choices"`
}

func ParseLine(line string) (*ChatResponse, bool) {
	line = strings.TrimSpace(line)
	if line == "" || line == "data: [DONE]" {
		return nil, false
	}
	line = strings.TrimPrefix(line, "data: ")
	var chunk ChatResponse
	if err := json.Unmarshal([]byte(line), &chunk); err != nil {
		return nil, false
	}
	return &chunk, true
}

func ParseLineContent(line string) string {
	if chunk, ok := ParseLine(line); ok && len(chunk.Choices) > 0 {
		c := chunk.Choices[0]
		if c.Delta.Content != "" {
			return c.Delta.Content
		}
		return c.Message.Content
	}
	return ""
}
