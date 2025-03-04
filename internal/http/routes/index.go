package routes

import (
	"VEDA95/open_board/api/internal/http/responses"
	"github.com/gofiber/fiber/v2"
)

func IndexGET(context *fiber.Ctx) error {
	return responses.JSONResponse(
		context,
		fiber.StatusOK,
		responses.OKResponse(fiber.StatusOK, responses.GenericMessage{Message: "Hello World"}),
	)
}
