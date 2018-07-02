package message

import (
	"encoding/json"
	"log"
)

type simpleMessage struct{
	Code int `json:"code"`
	Message string `json:"message"`
}

func BuildSimpleMessage(code int, msg string) []byte {
	var sm simpleMessage
	sm.Code = code
	sm.Message = msg
	b, err := json.Marshal(&sm)
	if err != nil {
		return nil
	}
	log.Println(string(b))
	return b
}