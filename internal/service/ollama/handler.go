package ollama

import (
	"ai-service/internal/repository"
	"ai-service/internal/util/config"
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const (
	chat          = "/api/chat"
	embed         = "/api/embed"
	role          = "system"
	system_prompt = "Вот текст документа, который ты должен использовать для ответа: "
)

type LLMService interface {
	Answer(ctx context.Context, orgID string, docEmbeds []float32, messages []Message) (*ChatResponse, error)
	Embed(input string) ([][]float32, error)
}

type llmService struct {
	config     *config.Config
	repository *repository.Repository
	client     *http.Client
}

func NewLLMService(cfg *config.Config, repo *repository.Repository) (LLMService, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{Transport: tr}
	httpClient.Timeout = cfg.Ollama.Timeout

	return &llmService{
		config:     cfg,
		repository: repo,
		client:     httpClient,
	}, nil
}

func (l *llmService) handler(method string, url string, req any) ([]byte, int, error) {
	body, err := json.Marshal(req)

	request, err := http.NewRequest(method, url, bytes.NewReader(body))
	request.Header.Add("Content-Type", "application/json")

	res, err := l.client.Do(request)
	if err != nil {
		return nil, 0, err
	}

	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, 0, err
	}

	return resBody, res.StatusCode, nil
}
