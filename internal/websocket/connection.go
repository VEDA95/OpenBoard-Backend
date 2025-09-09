package websocket

import (
	models "VEDA95/open_board/api/internal/db/model"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofrs/uuid/v5"
)

type AuthData struct {
	User         *models.User
	AccessToken  string
	Autheticated bool
}

type WebsocketConnection struct {
	ID          uuid.UUID
	Conncection *websocket.Conn
	AuthData    *AuthData
	Topics      *WebsocketTopicStore
}

func NewConnection(connection *websocket.Conn, user *models.User, accessToken string) (*WebsocketConnection, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	return &WebsocketConnection{
		ID:          id,
		Conncection: connection,
		Topics:      NewTopicStore(),
		AuthData: &AuthData{
			User:         user,
			AccessToken:  accessToken,
			Autheticated: true,
		},
	}, nil
}
