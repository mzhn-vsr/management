package handlers

import (
	"log/slog"
	"mzhn/management/internal/lib/logger/sl"
	"mzhn/management/internal/services/feedbackservice"

	"github.com/labstack/echo/v4"
)

// @Summary	Получение общей статистики по фидбеку
//
// @Tags		feedback
// @Success	200	{object} dto.FeedbackStats
// @Failure	500	{object}	InternalError
// @Router		/feedback [get]
func FeedbackStats(svc *feedbackservice.FeedbackService) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		stats, err := svc.Stats(ctx)
		if err != nil {
			slog.Error("cannot get stats", sl.Err(err))
			return c.JSON(500, &payload{
				"error": err.Error(),
			})
		}

		return c.JSON(200, stats)
	}
}
