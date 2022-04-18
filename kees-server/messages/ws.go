package messages

type WebSocket struct {
	State   string                 `json:"state"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}
