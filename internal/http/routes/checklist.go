package routes

import (
	"fmt"

	models "VEDA95/open_board/api/internal/db/model"
	"VEDA95/open_board/api/internal/errors"
	"VEDA95/open_board/api/internal/http/responses"
	"VEDA95/open_board/api/internal/http/validators"
	"VEDA95/open_board/api/internal/service"

	"github.com/gofiber/fiber/v2"
)

type CheckListItemHandler struct {
	checkListItemService *service.CheckListItemService
	validator            *validators.Validator
}

func NewCheckListItemHandler(checkListItemService *service.CheckListItemService, validator *validators.Validator) *CheckListItemHandler {
	return &CheckListItemHandler{
		checkListItemService: checkListItemService,
		validator:            validator,
	}
}

func (handler *CheckListItemHandler) GET(context *fiber.Ctx) error {
	cardID := context.Query("card_id")

	if cardID != "" {
		items, err := handler.checkListItemService.GetCheckListItemsByCardID(cardID)
		if err != nil {
			return err
		}
		return responses.JSONResponse(context, fiber.StatusOK, responses.OKCollectionResponse(fiber.StatusOK, items))
	}

	items, err := handler.checkListItemService.GetCheckListItems()
	if err != nil {
		return err
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKCollectionResponse(fiber.StatusOK, items))
}

func (handler *CheckListItemHandler) GETByID(context *fiber.Ctx) error {
	params := new(validators.ParamValidator)

	if err := context.ParamsParser(params); err != nil {
		return err
	}

	if errs := handler.validator.Validate(params); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	item, err := handler.checkListItemService.GetCheckListItemByID(params.Id)
	if err != nil {
		return err
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKResponse(fiber.StatusOK, item))
}

func (handler *CheckListItemHandler) POST(context *fiber.Ctx) error {
	validatorData := new(validators.CreateCheckListItemValidator)

	if err := context.BodyParser(validatorData); err != nil {
		return err
	}

	session := context.Locals("auth_session").(models.Session)
	validatorData.UserID = session.User.ID

	if errs := handler.validator.Validate(validatorData); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	item, err := handler.checkListItemService.CreateCheckListItem(validatorData)
	if err != nil {
		return err
	}

	return responses.JSONResponse(context, fiber.StatusCreated, responses.OKResponse(fiber.StatusCreated, item))
}

func (handler *CheckListItemHandler) PATCH(context *fiber.Ctx) error {
	params := new(validators.ParamValidator)

	if err := context.ParamsParser(params); err != nil {
		return err
	}

	if errs := handler.validator.Validate(params); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	validatorData := new(validators.UpdateCheckListItemValidator)

	if err := context.BodyParser(validatorData); err != nil {
		return err
	}

	if errs := handler.validator.Validate(validatorData); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	session := context.Locals("auth_session").(models.Session)

	item, err := handler.checkListItemService.UpdateCheckListItem(params.Id, session.User.ID, validatorData)
	if err != nil {
		return err
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKResponse(fiber.StatusOK, item))
}

func (handler *CheckListItemHandler) DELETE(context *fiber.Ctx) error {
	params := new(validators.ParamValidator)

	if err := context.ParamsParser(params); err != nil {
		return err
	}

	if errs := handler.validator.Validate(params); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	session := context.Locals("auth_session").(models.Session)

	if err := handler.checkListItemService.DeleteCheckListItem(params.Id, session.User.ID); err != nil {
		return err
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKResponse(fiber.StatusOK, responses.GenericMessage{
		Message: fmt.Sprintf("Checklist item %s has been deleted successfully", params.Id),
	}))
}

func (handler *CheckListItemHandler) Reorder(context *fiber.Ctx) error {
	params := validators.ParamValidator{Id: context.Params("card_id")}

	if errs := handler.validator.Validate(params); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	validatorData := new(validators.ReorderCheckListItemsValidator)

	if err := context.BodyParser(validatorData); err != nil {
		return err
	}

	if errs := handler.validator.Validate(validatorData); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	session := context.Locals("auth_session").(models.Session)

	if err := handler.checkListItemService.ReorderCheckListItems(session.User.ID, validatorData); err != nil {
		return err
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKResponse(fiber.StatusOK, responses.GenericMessage{
		Message: "Checklist items reordered successfully",
	}))
}
