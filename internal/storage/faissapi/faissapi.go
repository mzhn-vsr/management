package faissapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"mzhn/management/internal/config"
	"mzhn/management/internal/dto"
	"mzhn/management/internal/lib/logger/sl"
	"net/http"

	"github.com/labstack/gommon/log"
)

type FaissApi struct {
	url string
}

func New(cfg *config.FaissApi) *FaissApi {
	return &FaissApi{
		url: cfg.Url,
	}
}

func (api *FaissApi) Save(ctx context.Context, e []*dto.FaqFaissCreate) error {

	log.Debug("Marshaling json request", slog.Any("req", e))
	body, err := json.Marshal(e)
	if err != nil {
		log.Error("cannot marshal", sl.Err(err), slog.Any("req", e))
		return err
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/faiss/add", api.url), bytes.NewReader(body))
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	log.Debug("Response recieved", slog.Any("status", resp.StatusCode))
	if resp.StatusCode != 200 {
		log.Error("cannot add to faiss", slog.Int("status code", resp.StatusCode))
		return fmt.Errorf("cannot add to faiss")
	}

	return nil
}

func (api *FaissApi) Delete(ctx context.Context, ids []string) error {

	log.Debug("Marshaling json request", slog.Any("req", ids))
	body, err := json.Marshal(ids)
	if err != nil {
		log.Error("cannot marshal", sl.Err(err), slog.Any("req", ids))
		return err
	}

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/faiss/delete", api.url), bytes.NewReader(body))
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	log.Debug("Response recieved", slog.Any("status", resp.StatusCode))
	if resp.StatusCode != 200 {
		log.Error("cannot add to faiss", slog.Int("status code", resp.StatusCode))
		return fmt.Errorf("cannot add to faiss")
	}

	return nil
}
