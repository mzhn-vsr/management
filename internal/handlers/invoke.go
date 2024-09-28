package handlers

import (
	"mzhn/management/internal/services/chatservice"

	"github.com/labstack/echo/v4"
)

type PredictReq struct {
	Input string `json:"question"`
}

type PredictRes struct {
	A  string `json:"answer"`
	C1 string `json:"class_1"`
	C2 string `json:"class_2"`
}

// @Summary	Предикт ответа на вопрос
//
// @Param		input	body PredictReq	true	"input body"
// @Tags		feedback
// @Success	200	{object} PredictRes
// @Failure	500	{object}	InternalError
// @Router		/predict [post]
func Predict(cs *chatservice.ChatService) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req PredictReq

		if err := c.Bind(&req); err != nil {
			return c.JSON(500, &payload{
				"error": "cannot parse json",
			})
		}

		ctx := c.Request().Context()
		output, err := cs.Invoke(ctx, req.Input)
		if err != nil {
			return c.JSON(500, &payload{
				"error": err.Error(),
			})
		}

		return c.JSON(200, &PredictRes{
			A:  output.Answer,
			C1: output.Class1,
			C2: output.Class2,
		})
	}
}
