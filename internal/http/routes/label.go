package routes

import (
	"fmt"

	"VEDA95/open_board/api/internal/errors"
	"VEDA95/open_board/api/internal/http/responses"
	"VEDA95/open_board/api/internal/http/validators"
	"VEDA95/open_board/api/internal/service"

	"github.com/gofiber/fiber/v2"
)

type LabelHandler struct {
	labelService *service.LabelService
	validator    *validators.Validator
}

func NewLabelHandler(labelService *service.LabelService, validator *validators.Validator) *LabelHandler {
	return &LabelHandler{
		labelService: labelService,
		validator:    validator,
	}
}

func (handler *LabelHandler) GET(context *fiber.Ctx) error {
	boardID := context.Query("board_id")

	if boardID != "" {
		labels, err := handler.labelService.GetLabelsByBoardID(boardID)
		if err != nil {
			return err
		}
		return responses.JSONResponse(context, fiber.StatusOK, responses.OKCollectionResponse(fiber.StatusOK, labels))
	}

	labels, err := handler.labelService.GetLabels()
	if err != nil {
		return err
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKCollectionResponse(fiber.StatusOK, labels))
}

func (handler *LabelHandler) GETByID(context *fiber.Ctx) error {
	params := new(validators.ParamValidator)

	if err := context.ParamsParser(params); err != nil {
		return err
	}

	if errs := handler.validator.Validate(params); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	label, err := handler.labelService.GetLabelByID(params.Id)
	if err != nil {
		return err
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKResponse(fiber.StatusOK, label))
}

func (handler *LabelHandler) POST(context *fiber.Ctx) error {
	validatorData := new(validators.CreateLabelValidator)

	if err := context.BodyParser(validatorData); err != nil {
		return err
	}

	if errs := handler.validator.Validate(validatorData); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	label, err := handler.labelService.CreateLabel(validatorData)
	if err != nil {
		return err
	}

	return responses.JSONResponse(context, fiber.StatusCreated, responses.OKResponse(fiber.StatusCreated, label))
}

func (handler *LabelHandler) PATCH(context *fiber.Ctx) error {
	params := new(validators.ParamValidator)

	if err := context.ParamsParser(params); err != nil {
		return err
	}

	if errs := handler.validator.Validate(params); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	validatorData := new(validators.UpdateLabelValidator)

	if err := context.BodyParser(validatorData); err != nil {
		return err
	}

	if errs := handler.validator.Validate(validatorData); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	label, err := handler.labelService.UpdateLabel(params.Id, validatorData)
	if err != nil {
		return err
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKResponse(fiber.StatusOK, label))
}

func (handler *LabelHandler) DELETE(context *fiber.Ctx) error {
	params := new(validators.ParamValidator)

	if err := context.ParamsParser(params); err != nil {
		return err
	}

	if errs := handler.validator.Validate(params); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	if err := handler.labelService.DeleteLabel(params.Id); err != nil {
		return err
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKResponse(fiber.StatusOK, responses.GenericMessage{
		Message: fmt.Sprintf("Label %s has been deleted successfully", params.Id),
	}))
}
