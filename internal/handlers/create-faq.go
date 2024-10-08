package handlers

import (
	"log/slog"
	"mzhn/management/internal/dto"
	"mzhn/management/internal/lib/logger/sl"
	"mzhn/management/internal/services/faqservice"

	"github.com/labstack/echo/v4"
)

type CreateFaqRes struct {
	Id string `json:"id"`
}

// Создание записи в базе знаний
//
//	@Summary	Создание записи в БЗ
//
//	@Param		input	body	dto.FaqEntryCreate	true	"input body"
//
//	@Tags		faq
//	@Success	200	{object}	CreateFaqRes
//	@Failure	500	{object}	InternalError
//	@Router		/faq [post]
func CreateFaq(faqsvc *faqservice.FaqService) echo.HandlerFunc {
	return func(c echo.Context) error {
		log := slog.With(slog.String("handler", "CreateFaq"))
		var req dto.FaqEntryCreate

		if err := c.Bind(&req); err != nil {
			log.Error("cannot bind request", sl.Err(err))
			return c.JSON(500, &payload{
				"error": err.Error(),
			})
		}

		log.Debug("create for dto", slog.Any("req", req))

		ctx := c.Request().Context()
		id, err := faqsvc.Create(ctx, &req)
		if err != nil {
			log.Error("cannot create faq entry", sl.Err(err))
			return c.JSON(500, &payload{
				"error": err.Error(),
			})
		}

		return c.JSON(200, &payload{
			"id": id,
		})
	}
}
