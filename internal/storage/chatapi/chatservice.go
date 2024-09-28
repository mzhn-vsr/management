package chatapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mzhn/management/internal/config"
	"mzhn/management/internal/entity"
	"mzhn/management/internal/lib/logger/sl"
	"net/http"

	"github.com/labstack/gommon/log"
)

type ChatService struct {
	url string
	log *slog.Logger
}

func New(cfg *config.ChatApi) *ChatService {
	return &ChatService{
		log: slog.Default().With(slog.String("struct", "ChatService")),
		url: cfg.Url,
	}
}

func (c *ChatService) Invoke(ctx context.Context, input string) (*entity.ChatInvokeOutput, error) {

	var req struct {
		Input string `json:"input"`
	}

	req.Input = input

	log.Debug("Marshaling json request", slog.Any("req", req))
	body, err := json.Marshal(req)
	if err != nil {
		log.Error("cannot marshal", sl.Err(err), slog.Any("req", req))
		return nil, err
	}

	resp, err := http.Post(fmt.Sprintf("%s/chat/invoke", c.url), "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	log.Debug("Response recieved", slog.Any("status", resp.StatusCode))

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Error("error with reading response body", sl.Err(err))
		return nil, err
	}

	res := new(entity.ChatInvokeResponse)
	if err := json.Unmarshal(body, res); err != nil {
		log.Error("error with unmarshaling response body")
		return nil, err
	}

	return &entity.ChatInvokeOutput{
		Answer: res.Output.Content,
	}, nil
}
