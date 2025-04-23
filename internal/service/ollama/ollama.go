package ollama

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
)

func (l *llmService) Answer(ctx context.Context, orgID string, docEmbeds []float32, messages []Message) (*ChatResponse, error) {
	searchResult, err := l.repository.Vector.GetTopK(ctx, orgID, 5, docEmbeds)
	if err != nil {
		return nil, err
	}
	var documents string
	for _, result := range searchResult {
		column := result.Fields.GetColumn("text")
		for i := 0; i < column.Len(); i++ {
			s, err := column.GetAsString(i)
			if err != nil {
				return nil, err
			}
			documents += s
		}
	}
	messages = append(messages, Message{
		Role:    role,
		Content: system_prompt + documents,
	})
	messages[len(messages)-1], messages[len(messages)-2] = messages[len(messages)-2], messages[len(messages)-1]
	chatRequest := ChatRequest{
		Model:    l.config.Ollama.Model,
		Messages: messages,
		Stream:   false,
	}
	url := l.config.Ollama.Url + l.config.Ollama.Endpoints[chat]
	respBody, status, err := l.handler(http.MethodPost, url, chatRequest)
	if err != nil {
		return nil, err
	}

	switch status {
	case 200:
		var chatResponse ChatResponse
		err = json.Unmarshal(respBody, &chatResponse)
		if err != nil {
			return nil, err
		}
		return &chatResponse, nil
	default:
		return nil, errors.New(string(respBody))
	}
}

func (l *llmService) Embed(input []string) ([][][]float32, error) {
	var result [][][]float32
	for _, text := range input {
		req := EmbedRequest{
			Model: l.config.Ollama.Model,
			Input: text,
		}
		url := l.config.Ollama.Url + l.config.Ollama.Endpoints[embed]
		respBody, status, err := l.handler(http.MethodPost, url, req)
		if err != nil {
			return nil, err
		}

		switch status {
		case 200:
			var chatResponse EmbeddingResponse
			err = json.Unmarshal(respBody, &chatResponse)
			if err != nil {
				return nil, err
			}
			result = append(result, chatResponse.Embeddings)
		default:
			return nil, errors.New(string(respBody))
		}
	}
	return result, nil
}
