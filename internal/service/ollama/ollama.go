package ollama

import (
	"ai-service/internal/repository/milvus"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net/http"
)

func (l *llmService) Answer(docEmbeds []float32, messages []Message) (*ChatResponse, error) {
	searchResult, err := milvus.GetTopK(*l.milvus, 5, docEmbeds)
	if err != nil {
		return nil, err
	}
	var documents string
	for _, result := range searchResult {
		column := result.Fields.GetColumn("text")
		fmt.Println(column)
		for i := 0; i < column.Len(); i++ {
			s, err := column.GetAsString(i)
			if err != nil {
				return nil, err
			}
			documents += s
		}
	}
	messages = append(messages, Message{
		Role: "system",
		Content: "Вот текст документа, который ты должен использовать для ответа:" +
			documents,
	})
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

func (l *llmService) Embed(input string) ([][]float64, error) {
	req := EmbedRequest{
		Model: l.config.Ollama.Model,
		Input: input,
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
		return chatResponse.Embeddings, nil
	default:
		return nil, errors.New(string(respBody))
	}
}
