package validators

type WebsocketTopicSubriptionValidator struct {
	Type  string `json:"type" validator:"required,oneof=subscribe unsubscribe"`
	Topic string `json:"topic" validator:"required,min=1"`
}

type WebsocketAuthValidator struct {
	AccessToken string `json:"access_token" validator:"required,min=1"`
}
