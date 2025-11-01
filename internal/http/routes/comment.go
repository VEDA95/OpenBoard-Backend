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

type CommentHandler struct {
	commentService *service.CommentService
	validator      *validators.Validator
}

func NewCommentHandler(commentService *service.CommentService, validator *validators.Validator) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
		validator:      validator,
	}
}

func (handler *CommentHandler) GET(context *fiber.Ctx) error {
	cardID := context.Query("card_id")

	if cardID != "" {
		comments, err := handler.commentService.GetCommentsByCardID(cardID)
		if err != nil {
			return err
		}
		return responses.JSONResponse(context, fiber.StatusOK, responses.OKCollectionResponse(fiber.StatusOK, comments))
	}

	comments, err := handler.commentService.GetComments()
	if err != nil {
		return err
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKCollectionResponse(fiber.StatusOK, comments))
}

func (handler *CommentHandler) GETByID(context *fiber.Ctx) error {
	params := new(validators.ParamValidator)

	if err := context.ParamsParser(params); err != nil {
		return err
	}

	if errs := handler.validator.Validate(params); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	comment, err := handler.commentService.GetCommentByID(params.Id)
	if err != nil {
		return err
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKResponse(fiber.StatusOK, comment))
}

func (handler *CommentHandler) POST(context *fiber.Ctx) error {
	validatorData := new(validators.CreateCommentValidator)

	if err := context.BodyParser(validatorData); err != nil {
		return err
	}

	if errs := handler.validator.Validate(validatorData); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	session := context.Locals("auth_session").(models.Session)

	comment, err := handler.commentService.CreateComment(validatorData, session.User.ID)
	if err != nil {
		return err
	}

	return responses.JSONResponse(context, fiber.StatusCreated, responses.OKResponse(fiber.StatusCreated, comment))
}

func (handler *CommentHandler) PATCH(context *fiber.Ctx) error {
	params := new(validators.ParamValidator)

	if err := context.ParamsParser(params); err != nil {
		return err
	}

	if errs := handler.validator.Validate(params); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	validatorData := new(validators.UpdateCommentValidator)

	if err := context.BodyParser(validatorData); err != nil {
		return err
	}

	if errs := handler.validator.Validate(validatorData); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	session := context.Locals("auth_session").(models.Session)

	comment, err := handler.commentService.UpdateComment(params.Id, validatorData, session.User.ID)
	if err != nil {
		return err
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKResponse(fiber.StatusOK, comment))
}

func (handler *CommentHandler) DELETE(context *fiber.Ctx) error {
	params := new(validators.ParamValidator)

	if err := context.ParamsParser(params); err != nil {
		return err
	}

	if errs := handler.validator.Validate(params); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	session := context.Locals("auth_session").(models.Session)

	if err := handler.commentService.DeleteComment(params.Id, session.User.ID); err != nil {
		return err
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKResponse(fiber.StatusOK, responses.GenericMessage{
		Message: fmt.Sprintf("Comment %s has been deleted successfully", params.Id),
	}))
}
