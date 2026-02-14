package routes

import (
	"fmt"

	models "VEDA95/open_board/api/internal/db/model"
	"VEDA95/open_board/api/internal/errors"
	"VEDA95/open_board/api/internal/http/responses"
	"VEDA95/open_board/api/internal/http/validators"
	"VEDA95/open_board/api/internal/service"
	"VEDA95/open_board/api/internal/websocket"

	"github.com/gofiber/fiber/v2"
)

type CardHandler struct {
	cardService *service.CardService
	wsManager   *websocket.WebsocketConnectionManager
	validator   *validators.Validator
}

func NewCardHandler(cardService *service.CardService, wsManager *websocket.WebsocketConnectionManager, validator *validators.Validator) *CardHandler {
	return &CardHandler{
		cardService: cardService,
		wsManager:   wsManager,
		validator:   validator,
	}
}

func (cardHandler *CardHandler) GET(context *fiber.Ctx) error {
	cards, err := cardHandler.cardService.GetCards()
	if err != nil {
		return err
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKCollectionResponse(fiber.StatusOK, cards))
}

func (cardHandler *CardHandler) GETByID(context *fiber.Ctx) error {
	params := new(validators.ParamValidator)

	if err := context.ParamsParser(params); err != nil {
		return err
	}

	if errs := cardHandler.validator.Validate(params); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	card, err := cardHandler.cardService.GetCardByID(params.Id)
	if err != nil {
		return err
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKResponse(fiber.StatusOK, card))
}

func (cardHandler *CardHandler) POST(context *fiber.Ctx) error {
	validatorData := new(validators.CreateCardValidator)

	if err := context.BodyParser(validatorData); err != nil {
		return err
	}

	if errs := cardHandler.validator.Validate(validatorData); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	card, err := cardHandler.cardService.CreateCard(validatorData)
	if err != nil {
		return err
	}

	// Emit WebSocket event
	if cardHandler.wsManager != nil && card.List != nil && card.List.BoardID != "" {
		session := context.Locals("auth_session").(models.Session)
		cardHandler.wsManager.EmitToTopic(
			session.User.ID,
			websocket.BoardTopic(card.List.BoardID),
			websocket.WebsocketMessage{
				Type: websocket.EventCardCreated,
				Data: map[string]any{
					"card_id":  card.ID,
					"board_id": card.List.BoardID,
					"card":     card,
				},
			},
		)
	}

	return responses.JSONResponse(context, fiber.StatusCreated, responses.OKResponse(fiber.StatusCreated, card))
}

func (cardHandler *CardHandler) PATCH(context *fiber.Ctx) error {
	params := new(validators.ParamValidator)

	if err := context.ParamsParser(params); err != nil {
		return err
	}

	if errs := cardHandler.validator.Validate(params); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	validatorData := new(validators.UpdateCardValidator)

	if err := context.BodyParser(validatorData); err != nil {
		return err
	}

	if errs := cardHandler.validator.Validate(validatorData); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	card, err := cardHandler.cardService.UpdateCard(params.Id, validatorData)
	if err != nil {
		return err
	}

	// Emit WebSocket event
	if cardHandler.wsManager != nil && card.List != nil && card.List.BoardID != "" {
		session := context.Locals("auth_session").(models.Session)
		cardHandler.wsManager.EmitToTopic(
			session.User.ID,
			websocket.BoardTopic(card.List.BoardID),
			websocket.WebsocketMessage{
				Type: websocket.EventCardUpdated,
				Data: map[string]any{
					"card_id":  card.ID,
					"board_id": card.List.BoardID,
					"card":     card,
				},
			},
		)
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKResponse(fiber.StatusOK, card))
}

func (cardHandler *CardHandler) Move(context *fiber.Ctx) error {
	params := new(validators.ParamValidator)

	if err := context.ParamsParser(params); err != nil {
		return err
	}

	if errs := cardHandler.validator.Validate(params); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	validatorData := new(validators.MoveCardValidator)

	if err := context.BodyParser(validatorData); err != nil {
		return err
	}

	if errs := cardHandler.validator.Validate(validatorData); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	card, err := cardHandler.cardService.MoveCard(params.Id, validatorData)
	if err != nil {
		return err
	}

	// Emit WebSocket event
	if cardHandler.wsManager != nil && card.List != nil && card.List.BoardID != "" {
		session := context.Locals("auth_session").(models.Session)
		cardHandler.wsManager.EmitToTopic(
			session.User.ID,
			websocket.BoardTopic(card.List.BoardID),
			websocket.WebsocketMessage{
				Type: websocket.EventCardMoved,
				Data: map[string]any{
					"card_id":  card.ID,
					"board_id": card.List.BoardID,
					"card":     card,
				},
			},
		)
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKResponse(fiber.StatusOK, card))
}

func (cardHandler *CardHandler) DELETE(context *fiber.Ctx) error {
	params := new(validators.ParamValidator)

	if err := context.ParamsParser(params); err != nil {
		return err
	}

	if errs := cardHandler.validator.Validate(params); len(errs) > 0 {
		return errors.CreateValidationError(errs)
	}

	// Get the card first to find the board ID
	card, err := cardHandler.cardService.GetCardByID(params.Id)
	if err != nil {
		return err
	}

	boardID := ""
	if card.List != nil {
		boardID = card.List.BoardID
	}

	if err := cardHandler.cardService.DeleteCard(params.Id); err != nil {
		return err
	}

	// Emit WebSocket event
	if cardHandler.wsManager != nil && boardID != "" {
		session := context.Locals("auth_session").(models.Session)
		cardHandler.wsManager.EmitToTopic(
			session.User.ID,
			websocket.BoardTopic(boardID),
			websocket.WebsocketMessage{
				Type: websocket.EventCardDeleted,
				Data: map[string]any{
					"card_id":  params.Id,
					"board_id": boardID,
				},
			},
		)
	}

	return responses.JSONResponse(context, fiber.StatusOK, responses.OKResponse(fiber.StatusOK, responses.GenericMessage{
		Message: fmt.Sprintf("Card %s has been deleted successfully", params.Id),
	}))
}

func (cardHandler *CardHandler) AddComment(context *fiber.Ctx) error {
	return fiber.NewError(fiber.StatusNotImplemented, "Comment functionality not yet implemented")
}
