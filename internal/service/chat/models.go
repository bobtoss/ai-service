package chat

import "ai-service/internal/service/ollama"

type ChatRequest struct {
	RqUID    string `json:"rquid"`
	Messages []ollama.Message
}

type ChatResponse struct {
	Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
}
