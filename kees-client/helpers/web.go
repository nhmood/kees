package helpers

import (
	"encoding/json"
)

func Format(data interface{}) ([]byte, error) {
	output, err := json.Marshal(data)
	return output, err
}

func ToInterface(data interface{}) map[string]interface{} {
	str, _ := Format(data)
	var conv map[string]interface{}
	json.Unmarshal([]byte(str), &conv)
	return conv
}

func ToStruct(data interface{}, target interface{}) {
	str, _ := Format(data)
	json.Unmarshal([]byte(str), target)
}
