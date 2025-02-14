package ollama

import (
	"ai-service/internal/util/config"
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"io/ioutil"
	"net/http"
)

const (
	chat  = "/api/chat"
	embed = "/api/embed"
)

type LLMService interface {
	Answer(docEmbeds []float32, messages []Message) (*ChatResponse, error)
	Embed(input string) ([][]float64, error)
}

type llmService struct {
	config *config.Config
	milvus *client.Client
	client *http.Client
}

func NewLLMService(ctx context.Context, cfg *config.Config) (LLMService, error) {
	milvus, err := client.NewClient(ctx, client.Config{Address: cfg.Milvus.Host})
	if err != nil {
		return nil, err
	}
	defer milvus.Close()

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{Transport: tr}
	httpClient.Timeout = cfg.Ollama.Timeout

	return &llmService{
		config: cfg,
		milvus: &milvus,
		client: httpClient,
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
