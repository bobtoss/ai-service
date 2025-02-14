package doc

import (
	"ai-service/internal/repository"
	"ai-service/internal/service/ollama"
	"ai-service/internal/util/config"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"github.com/labstack/echo"
	"io/ioutil"
	"net/http"
)

const (
	chat  = "/api/chat"
	embed = "/api/embed"
)

type DocService interface {
	SaveDoc(c echo.Context) error
	DeleteDoc(c echo.Context) error
	UpdatePriority(c echo.Context) error
}

type docService struct {
	config *config.Config
	llm    ollama.LLMService
	client *http.Client
}

func NewDocService(cfg *config.Config, repo *repository.Repository) (DocService, error) {
	llm, err := ollama.NewLLMService(cfg, repo)
	if err != nil {
		return nil, err
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{Transport: tr}
	httpClient.Timeout = cfg.Ollama.Timeout

	return &docService{
		config: cfg,
		llm:    llm,
		client: httpClient,
	}, nil
}

func (l *docService) handler(method string, url string, req any) ([]byte, int, error) {
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
