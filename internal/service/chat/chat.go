package chat

import (
	"ai-service/internal/repository"
	"ai-service/internal/service/ollama"
	"ai-service/internal/util/config"
	"ai-service/internal/util/errors"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
)

const (
	chat  = "/api/chat"
	embed = "/api/embed"
)

type ChatService interface {
	Chat(c echo.Context) error
}

type chatService struct {
	config     *config.Config
	llm        ollama.LLMService
	repository *repository.Repository
	client     *http.Client
}

func NewChatService(cfg *config.Config, repo *repository.Repository) (ChatService, error) {
	llm, err := ollama.NewLLMService(cfg, repo)
	if err != nil {
		return nil, err
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{Transport: tr}
	httpClient.Timeout = cfg.Ollama.Timeout

	return &chatService{
		config:     cfg,
		repository: repo,
		llm:        llm,
		client:     httpClient,
	}, nil
}

// Chat
//
// @Description Chat
// @Summary	Save document in milvus
// @Tags doc
// @Accept json
// @Produce	json
// @Success	200				{object}		models.SaveDocResponse
// @Failure	500				{object}	status.SaveDoc
// @Router /api/v1/chat 	[post]
func (d *chatService) Chat(c echo.Context) error {
	ctx := c.Request().Context()
	var dataReq ChatRequest
	if err := c.Bind(&dataReq); err != nil {
		return errors.NewBadRequestErrorRsp(err.Error())
	}

	embeddings, err := d.llm.Embed(SliceFromMesages(dataReq))
	if err != nil {
		return errors.NewInternalErrorRsp(err.Error())
	}
	response, err := d.llm.Answer(ctx, "new", embeddings[0][0], dataReq.Messages)
	if err != nil {
		return errors.NewInternalErrorRsp(err.Error())
	}
	return c.JSON(http.StatusOK, response)
}

func SliceFromMesages(messages ChatRequest) []string {
	var slice []string
	slice = append(slice, messages.Messages[len(messages.Messages)-1].Content)
	return slice
}

func (d *chatService) handler(method string, url string, req any) ([]byte, int, error) {
	body, err := json.Marshal(req)

	request, err := http.NewRequest(method, url, bytes.NewReader(body))
	request.Header.Add("Content-Type", "application/json")

	res, err := d.client.Do(request)
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
