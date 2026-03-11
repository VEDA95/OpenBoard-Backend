package routes

import (
	"VEDA95/open_board/api/internal/errors"
	"VEDA95/open_board/api/internal/http/responses"
	"VEDA95/open_board/api/internal/http/validators"
	"VEDA95/open_board/api/internal/service"
	"fmt"

	models "VEDA95/open_board/api/internal/db/model"

	"github.com/gofiber/fiber/v2"
)

type BoardHandler struct {
	boardService *service.BoardService
	validator    *validators.Validator
}

func NewBoardHandler(boardService *service.BoardService, validator *validators.Validator) *BoardHandler {
	return &BoardHandler{
		boardService: boardService,
		validator:    validator,
	}
}

func (boardHandler *BoardHandler) GET(context *fiber.Ctx) error {
	boards, err := boardHandler.boardService.GetBoards()
	if err != nil {
		return err
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKCollectionResponse(fiber.StatusOK, boards))
}

func (boardHandler *BoardHandler) GETByID(context *fiber.Ctx) error {
	params := new(validators.ParamValidator)

	if err := context.ParamsParser(params); err != nil {
		return err
	}

	if errs := boardHandler.validator.Validate(params); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	board, err := boardHandler.boardService.GetBoardByID(params.Id)
	if err != nil {
		return err
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKResponse(fiber.StatusOK, board))
}

func (boardHandler *BoardHandler) POST(context *fiber.Ctx) error {
	session := context.Locals("auth_session").(models.Session)
	validatorData := new(validators.CreateBoardValidator)

	if err := context.BodyParser(validatorData); err != nil {
		return err
	}

	if errs := boardHandler.validator.Validate(validatorData); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	board, err := boardHandler.boardService.CreateBoard(session.User, validatorData)
	if err != nil {
		return err
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKResponse(fiber.StatusOK, board))
}

func (boardHandler *BoardHandler) PATCH(context *fiber.Ctx) error {
	params := new(validators.ParamValidator)

	if err := context.ParamsParser(params); err != nil {
		return err
	}

	if errs := boardHandler.validator.Validate(params); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	validatorData := new(validators.UpdateBoardValidator)

	if err := context.BodyParser(validatorData); err != nil {
		return err
	}

	if errs := boardHandler.validator.Validate(validatorData); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	session := context.Locals("auth_session").(models.Session)
	board, err := boardHandler.boardService.UpdateBoard(params.Id, session.User, validatorData)
	if err != nil {
		return err
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKResponse(fiber.StatusOK, board))
}

func (boardHandler *BoardHandler) DELETE(context *fiber.Ctx) error {
	params := new(validators.ParamValidator)

	if err := context.ParamsParser(params); err != nil {
		return err
	}

	if errs := boardHandler.validator.Validate(params); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	if err := boardHandler.boardService.DeleteBoard(params.Id); err != nil {
		return err
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKResponse(fiber.StatusOK, responses.GenericMessage{Message: fmt.Sprintf("board: %s has been deleted successfully", params.Id)}))
}
