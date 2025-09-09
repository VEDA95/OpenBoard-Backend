package websocket

import (
	models "VEDA95/open_board/api/internal/db/model"
	"VEDA95/open_board/api/internal/http/validators"
	"VEDA95/open_board/api/internal/log"
	"VEDA95/open_board/api/internal/service"
	"fmt"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type WebsocketConnectionManager struct {
	Store       *WebsocketConnectionStore
	authService *service.AuthService
	validator   *validators.Validator
}

func NewWebsocketConnectionManager(authService *service.AuthService, validator *validators.Validator) *WebsocketConnectionManager {
	return &WebsocketConnectionManager{
		Store:       NewConnectionStore(),
		authService: authService,
		validator:   validator,
	}
}

func (connectionManager *WebsocketConnectionManager) checkUserSession(connection *WebsocketConnection) {
	ticker := time.NewTicker(time.Minute * 15)

	defer ticker.Stop()

	for range ticker.C {
		_, err := connectionManager.authService.ValidateSession(connection.AuthData.AccessToken)
		if err == nil {
			continue
		}

		connection.AuthData.Autheticated = false
		errMessage := err.Error()
		if err := connection.Conncection.WriteJSON(WebsocketMessage{Type: "auth", Message: &errMessage}); err != nil {
			log.Global.Error().Err(err).Msgf("an error occurred when sending message to websocket connection: %s", connection.ID.String())
		}
		break
	}
}

func (*WebsocketConnectionManager) RequireConnectionUpgrade(context *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(context) {
		return context.Next()
	}

	return fiber.ErrUpgradeRequired
}

func (connectionManager *WebsocketConnectionManager) ListenToConnection() fiber.Handler {
	return websocket.New(func(rawConnection *websocket.Conn) {
		session := rawConnection.Locals("auth_session").(models.Session)
		connection, err := NewConnection(rawConnection, session.User, session.AccessToken)
		if err != nil {
			log.Global.Error().Err(err).Msg("an error occurred when initializing websocket connection")
			return
		}

		connectionManager.Store.Set(connection)
		go connectionManager.checkUserSession(connection)

		for {
			if !connection.AuthData.Autheticated {
				message := new(validators.WebsocketAuthValidator)
				if err := connection.Conncection.ReadJSON(message); err != nil {
					log.Global.Error().Err(err).Msgf("an error occurred when reading message from websocket connection: %s", connection.ID.String())
					break
				}

				if errs := connectionManager.validator.Validate(message); len(errs) > 0 {
					go func() {
						errMessage := "one or more of the following fields are invalid"
						if err := connection.Conncection.WriteJSON(WebsocketMessage{Type: "auth", Message: &errMessage, Data: errs}); err != nil {
							log.Global.Error().Err(err).Msgf("an error occurred when sending message to websocket connection: %s", connection.ID.String())
						}
					}()
					continue
				}

				session, err := connectionManager.authService.ValidateSession(message.AccessToken)
				if err != nil {
					go func() {
						errMessage := "unauthorized"
						if err := connection.Conncection.WriteJSON(WebsocketMessage{Type: "auth", Message: &errMessage}); err != nil {
							log.Global.Error().Err(err).Msgf("an error occurred when sending message to websocket connection: %s", connection.ID.String())
						}
					}()
					continue
				}

				connection.AuthData.User = session.User
				connection.AuthData.AccessToken = session.AccessToken
				connection.AuthData.Autheticated = true

				go connectionManager.checkUserSession(connection)
				go func() {
					successMessage := "auth check successful!"
					if err := connection.Conncection.WriteJSON(WebsocketMessage{Type: "auth", Message: &successMessage}); err != nil {
						log.Global.Error().Err(err).Msgf("an error occurred when sending message to websocket connection: %s", connection.ID.String())
					}
				}()
				continue
			}

			message := new(validators.WebsocketTopicSubriptionValidator)
			if err := connection.Conncection.ReadJSON(message); err != nil {
				log.Global.Error().Err(err).Msgf("an error occurred when reading message from websocket connection: %s", connection.ID.String())
				break
			}

			if errs := connectionManager.validator.Validate(message); len(errs) > 0 {
				go func() {
					errMessage := "one or more of the following fields are invalid"
					if err := connection.Conncection.WriteJSON(WebsocketMessage{Type: "topic", Message: &errMessage, Data: errs}); err != nil {
						log.Global.Error().Err(err).Msgf("an error occurred when sending message to websocket connection: %s", connection.ID.String())
					}
				}()
				continue
			}

			if message.Type == "subscribe" {
				if connection.Topics.Has(message.Topic) {
					go func() {
						successMessage := "the provided topic has already been subscribed to"
						if err := connection.Conncection.WriteJSON(WebsocketMessage{Type: "subscribe", Message: &successMessage}); err != nil {
							log.Global.Error().Err(err).Msgf("an error occurred when sending message to websocket connection: %s", connection.ID.String())
						}
					}()
					continue
				}

				connection.Topics.Set(message.Topic)

				go func() {
					successMessage := fmt.Sprintf("topic: %s has been added successfully", message.Topic)
					if err := connection.Conncection.WriteJSON(WebsocketMessage{Type: "subscribe", Message: &successMessage}); err != nil {
						log.Global.Error().Err(err).Msgf("an error occurred when sending message to websocket connection: %s", connection.ID.String())
					}
				}()
			}

			if message.Type == "unsubscribe" {
				if !connection.Topics.Has(message.Topic) {
					go func() {
						successMessage := "the provided topic does not exist"
						if err := connection.Conncection.WriteJSON(WebsocketMessage{Type: "unsubscribe", Message: &successMessage}); err != nil {
							log.Global.Error().Err(err).Msgf("an error occurred when sending message to websocket connection: %s", connection.ID.String())
						}
					}()
					continue
				}

				connection.Topics.Remove(message.Topic)

				go func() {
					successMessage := fmt.Sprintf("topic: %s has been removed successfully", message.Topic)
					if err := connection.Conncection.WriteJSON(WebsocketMessage{Type: "unsubscribe", Message: &successMessage}); err != nil {
						log.Global.Error().Err(err).Msgf("an error occurred when sending message to websocket connection: %s", connection.ID.String())
					}
				}()
			}
		}

		connectionManager.Store.Remove(connection.ID.String())
	})
}

func (connectionManager *WebsocketConnectionManager) Send(ID string, message WebsocketMessage) error {
	connection, err := connectionManager.Store.Get(ID)
	if err != nil {
		return err
	}

	if err := connection.Conncection.WriteJSON(message); err != nil {
		return err
	}

	return nil
}

func (connectionManager *WebsocketConnectionManager) BroadcastAll(message WebsocketMessage) {
	connections := connectionManager.Store.GetAll()

	for _, connection := range connections {
		go func() {
			if err := connection.Conncection.WriteJSON(message); err != nil {
				log.Global.Error().Err(err).Msgf("an error occurred when sending message to websocket connection: %s", connection.ID.String())
			}
		}()
	}
}

func (connectionManager *WebsocketConnectionManager) BroadcastToTopic(topic string, message WebsocketMessage) {
	connections := connectionManager.Store.GetByTopic(topic)

	for _, connection := range connections {
		go func() {
			if err := connection.Conncection.WriteJSON(message); err != nil {
				log.Global.Error().Err(err).Msgf("an error occurred when sending message to websocket connection: %s", connection.ID.String())
			}
		}()
	}
}
