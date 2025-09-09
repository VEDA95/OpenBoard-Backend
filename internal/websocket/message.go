package websocket

type WebsocketMessage struct {
	Message *string `json:"message,omitempty"`
	Type    string  `json:"type"`
	Data    any     `json:"data,omitempty"`
}
