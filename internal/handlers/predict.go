package handlers

import (
	"mzhn/management/internal/services/chatservice"

	"github.com/labstack/echo/v4"
)

func Invoke(cs *chatservice.ChatService) echo.HandlerFunc {

	type request struct {
		Input string `json:"question"`
	}

	type response struct {
		FeedbackId string `json:"feedbackId"`
		Answer     string `json:"answer"`
	}

	return func(c echo.Context) error {
		var req request

		if err := c.Bind(&req); err != nil {
			return c.JSON(500, &payload{
				"error": "cannot parse json",
			})
		}

		ctx := c.Request().Context()
		output, err := cs.Predicate(ctx, req.Input)
		if err != nil {
			return c.JSON(500, &payload{
				"error": err.Error(),
			})
		}

		return c.JSON(200, &response{
			Answer:     output.Answer,
			FeedbackId: output.Id,
		})
	}
}
