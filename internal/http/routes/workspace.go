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

type WorkspaceHandler struct {
	workspaceService *service.WorkspaceService
	validator        *validators.Validator
}

func NewWorkspaceHandler(workspaceService *service.WorkspaceService, validator *validators.Validator) *WorkspaceHandler {
	return &WorkspaceHandler{
		workspaceService: workspaceService,
		validator:        validator,
	}
}

func (workspaceHandler *WorkspaceHandler) GET(context *fiber.Ctx) error {
	session := context.Locals("auth_session").(models.Session)
	workspaces, err := workspaceHandler.workspaceService.GetUserAccessibleWorkspaces(session.User)
	if err != nil {
		return err
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKCollectionResponse(fiber.StatusOK, workspaces))
}

func (workspaceHandler *WorkspaceHandler) GETByID(context *fiber.Ctx) error {
	params := new(validators.ParamValidator)

	if err := context.ParamsParser(params); err != nil {
		return err
	}

	if errs := workspaceHandler.validator.Validate(params); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	workspace, err := workspaceHandler.workspaceService.GetWorkspaceByID(params.Id)
	if err != nil {
		return err
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKResponse(fiber.StatusOK, workspace))
}

func (workspaceHandler *WorkspaceHandler) POST(context *fiber.Ctx) error {
	session := context.Locals("auth_session").(models.Session)
	validatorData := new(validators.CreateWorkspaceValidator)

	if err := context.BodyParser(validatorData); err != nil {
		return err
	}

	if errs := workspaceHandler.validator.Validate(validatorData); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	workspace, err := workspaceHandler.workspaceService.CreateWorkspace(session.User, validatorData)
	if err != nil {
		return err
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKResponse(fiber.StatusOK, workspace))
}

func (workspaceHandler *WorkspaceHandler) PATCH(context *fiber.Ctx) error {
	params := new(validators.ParamValidator)

	if err := context.ParamsParser(params); err != nil {
		return err
	}

	if errs := workspaceHandler.validator.Validate(params); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	validatorData := new(validators.UpdateWorkspaceValidator)

	if err := context.BodyParser(validatorData); err != nil {
		return err
	}

	if errs := workspaceHandler.validator.Validate(validatorData); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	session := context.Locals("auth_session").(models.Session)
	workspace, err := workspaceHandler.workspaceService.UpdateWorkspace(params.Id, session.User, validatorData)
	if err != nil {
		return err
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKResponse(fiber.StatusOK, workspace))
}

func (workspaceHandler *WorkspaceHandler) DELETE(context *fiber.Ctx) error {
	params := new(validators.ParamValidator)

	if err := context.ParamsParser(params); err != nil {
		return err
	}

	if errs := workspaceHandler.validator.Validate(params); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	if err := workspaceHandler.workspaceService.DeleteWorkspace(params.Id); err != nil {
		return err
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKResponse(fiber.StatusOK, responses.GenericMessage{Message: fmt.Sprintf("Workspace: %s has been deleted successfully", params.Id)}))
}
