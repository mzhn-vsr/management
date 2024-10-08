package handlers

import (
	"errors"
	"log/slog"
	"mzhn/management/internal/lib/logger/sl"
	"mzhn/management/internal/services/faqservice"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// @Summary	Поиск конкретной записи вопрос-ответ из базы знаний
//
// @Param		id	path	int	true	"faq id"
// @Tags		faq
// @Success	200	{object} dto.FeedbackStats
// @Failure	401	{object}	InternalError
// @Failure	404	{object}	InternalError
// @Failure	500	{object}	InternalError
// @Router		/faq/{id} [get]
func FindFaq(faqsvc *faqservice.FaqService) echo.HandlerFunc {
	return func(c echo.Context) error {
		log := slog.With(slog.String("handler", "UpdateFaq"))

		id := c.Param("id")
		if id == "" {
			return c.JSON(echo.ErrBadRequest.Code, &payload{
				"error": "empty id",
			})
		}
		if _, err := uuid.Parse(id); err != nil {
			return c.JSON(echo.ErrBadRequest.Code, &payload{
				"error": "id must be uuid",
			})
		}

		ctx := c.Request().Context()
		entry, err := faqsvc.Find(ctx, id)
		if err != nil {
			if errors.Is(err, faqservice.ErrNotFound) {
				return c.JSON(404, &payload{
					"error": "faq entry does not exists with id " + id,
				})
			}
			log.Error("cannot create faq entry", sl.Err(err))
			return c.JSON(500, &payload{
				"error": err.Error(),
			})
		}

		return c.JSON(200, entry)
	}
}
