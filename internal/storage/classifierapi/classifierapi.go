package classifierapi

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
)

type ClassifierApi struct {
	url    string
	logger *slog.Logger
}

func New(cfg *config.ClassifierApi) *ClassifierApi {
	return &ClassifierApi{
		url:    cfg.Url,
		logger: slog.Default().With(slog.String("struct", "ClassifierApi")),
	}
}

func (api *ClassifierApi) Classify(ctx context.Context, input string) (*entity.ClassifierResponse, error) {

	log := api.logger.With(slog.String("method", "Classify"))

	var req struct {
		Answer string `json:"answer"`
	}

	req.Answer = input

	log.Debug("Marshaling json request", slog.Any("req", req))
	body, err := json.Marshal(req)
	if err != nil {
		log.Error("cannot marshal", sl.Err(err), slog.Any("req", req))
		return nil, err
	}

	endpoint := fmt.Sprintf("%s/predict", api.url)
	log.Debug("sending post request", slog.String("endpoint", endpoint), slog.Any("req", req))
	resp, err := http.Post(endpoint, "application/json", bytes.NewReader(body))
	if err != nil {
		log.Error("cannot send post request", slog.String("endpoint", endpoint), sl.Err(err))
		return nil, err
	}
	defer resp.Body.Close()

	var res struct {
		Class1 string `json:"prediction"`
		Class2 string `json:"class_2"`
	}

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Debug("response recieved", slog.String("response", string(resBody)))

	if err := json.Unmarshal(resBody, &res); err != nil {
		return nil, err
	}
	log.Debug("response body unmarshal", slog.Any("body", res))

	return &entity.ClassifierResponse{
		Class1: res.Class1,
		Class2: res.Class2,
	}, nil
}
