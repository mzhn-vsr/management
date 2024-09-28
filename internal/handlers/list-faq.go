package handlers

import (
	"log/slog"
	"mzhn/management/internal/dto"
	"mzhn/management/internal/entity"
	"mzhn/management/internal/lib/logger/sl"
	"mzhn/management/internal/services/faqservice"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ListFaqRes struct {
	Items []entity.FaqEntry `json:"items"`
	Total uint64            `json:"total"`
}

// @Summary	Получение списка вопрос-ответ из базы знаний
//
// @Param		offset	query int	true	"input body"
// @Param		limit	query int	true	"input body"
// @Tags		faq
// @Success	200	{object} ListFaqRes
// @Failure	500	{object}	InternalError
// @Router		/faq [get]
func ListFaq(faqsvc *faqservice.FaqService) echo.HandlerFunc {
	return func(c echo.Context) error {
		log := slog.With(slog.String("handler", "ListFaq"))

		filters := dto.FaqEntryList{}

		offset := c.QueryParam("offset")
		if offset != "" {
			offset, err := strconv.ParseUint(offset, 10, 32)
			if err != nil {
				return c.JSON(echo.ErrBadRequest.Code, &payload{
					"error": "offset has invalid format",
				})
			}
			filters.Offset = &offset
		}

		limit := c.QueryParam("limit")
		if limit != "" {
			limit, err := strconv.ParseUint(limit, 10, 32)
			if err != nil {
				return c.JSON(echo.ErrBadRequest.Code, &payload{
					"error": "limit has invalid format",
				})
			}

			filters.Limit = &limit
		}

		ctx := c.Request().Context()
		entries, total, err := faqsvc.List(ctx, filters)
		if err != nil {
			log.Error("cannot list faq entry", sl.Err(err))
			return c.JSON(500, &payload{
				"error": err.Error(),
			})
		}

		return c.JSON(200, &payload{
			"items": entries,
			"total": total,
		})
	}
}
