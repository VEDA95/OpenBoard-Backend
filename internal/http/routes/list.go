package routes

import (
	"fmt"

	"VEDA95/open_board/api/internal/errors"
	"VEDA95/open_board/api/internal/http/responses"
	"VEDA95/open_board/api/internal/http/validators"
	"VEDA95/open_board/api/internal/service"

	"github.com/gofiber/fiber/v2"
)

type ListHandler struct {
	listService *service.ListService
	validator   *validators.Validator
}

func NewListHandler(listService *service.ListService, validator *validators.Validator) *ListHandler {
	return &ListHandler{
		listService: listService,
		validator:   validator,
	}
}

func (listHandler *ListHandler) GET(context *fiber.Ctx) error {
	boardID := context.Query("board_id")

	if boardID != "" {
		lists, err := listHandler.listService.GetListsByBoardID(boardID)
		if err != nil {
			return err
		}
		return responses.JSONResponse(context, fiber.StatusOK, responses.OKCollectionResponse(fiber.StatusOK, lists))
	}

	lists, err := listHandler.listService.GetLists()
	if err != nil {
		return err
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKCollectionResponse(fiber.StatusOK, lists))
}

func (listHandler *ListHandler) GETByID(context *fiber.Ctx) error {
	params := new(validators.ParamValidator)

	if err := context.ParamsParser(params); err != nil {
		return err
	}

	if errs := listHandler.validator.Validate(params); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	list, err := listHandler.listService.GetListByID(params.Id)
	if err != nil {
		return err
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKResponse(fiber.StatusOK, list))
}

func (listHandler *ListHandler) POST(context *fiber.Ctx) error {
	validatorData := new(validators.CreateListValidator)

	if err := context.BodyParser(validatorData); err != nil {
		return err
	}

	if errs := listHandler.validator.Validate(validatorData); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	list, err := listHandler.listService.CreateList(validatorData)
	if err != nil {
		return err
	}

	return responses.JSONResponse(context, fiber.StatusCreated, responses.OKResponse(fiber.StatusCreated, list))
}

func (listHandler *ListHandler) PATCH(context *fiber.Ctx) error {
	params := new(validators.ParamValidator)

	if err := context.ParamsParser(params); err != nil {
		return err
	}

	if errs := listHandler.validator.Validate(params); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	validatorData := new(validators.UpdateListValidator)

	if err := context.BodyParser(validatorData); err != nil {
		return err
	}

	if errs := listHandler.validator.Validate(validatorData); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	list, err := listHandler.listService.UpdateList(params.Id, validatorData)
	if err != nil {
		return err
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKResponse(fiber.StatusOK, list))
}

func (listHandler *ListHandler) DELETE(context *fiber.Ctx) error {
	params := new(validators.ParamValidator)

	if err := context.ParamsParser(params); err != nil {
		return err
	}

	if errs := listHandler.validator.Validate(params); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	if err := listHandler.listService.DeleteList(params.Id); err != nil {
		return err
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKResponse(fiber.StatusOK, responses.GenericMessage{
		Message: fmt.Sprintf("List %s has been deleted successfully", params.Id),
	}))
}
