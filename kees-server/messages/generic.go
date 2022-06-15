package messages

type Generic struct {
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}
