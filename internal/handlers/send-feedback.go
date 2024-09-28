package handlers

import (
	"errors"
	"fmt"
	"log/slog"
	"mzhn/management/internal/lib/logger/sl"
	"mzhn/management/internal/services/feedbackservice"

	"github.com/labstack/echo/v4"
)

func SendFeedback(sender *feedbackservice.FeedbackService) echo.HandlerFunc {
	type request struct {
		Id       string `json:"id"`
		IsUseful bool   `json:"isUseful"`
	}

	return func(c echo.Context) error {
		var req request

		if err := c.Bind(&req); err != nil {
			return c.JSON(500, &payload{
				"error": err.Error(),
			})
		}

		slog.Debug("recieved send feedback request", slog.Any("req", req))

		ctx := c.Request().Context()
		if err := sender.Send(ctx, req.Id, req.IsUseful); err != nil {

			if errors.Is(err, feedbackservice.ErrNoMessage) {
				slog.Debug("message not found", slog.String("id", req.Id))
				return c.JSON(404, &payload{
					"id":    req.Id,
					"error": fmt.Sprintf("message with id = %s not found", req.Id),
				})
			}

			if errors.Is(err, feedbackservice.ErrFeedbackAlreadySent) {
				slog.Debug("feedback already sent for message", slog.String("messageId", req.Id))
				return c.JSON(400, &payload{
					"error": "feedback already sent",
				})
			}

			slog.Debug("cannot send", sl.Err(err))
			return c.JSON(500, &payload{
				"error": err.Error(),
			})
		}

		return c.JSON(200, &payload{
			"message": "feedback sent",
		})
	}
}
