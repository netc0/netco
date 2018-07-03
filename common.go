package netco

import (
	"encoding/json"
	"log"
	"hash/crc32"
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

func BuildPushPacket(route string, data []byte) []byte{
	if data == nil { data = make([]byte, 0)}
	var result = make([]byte, 8)
	var bodyLen = 4 + uint32(len(data))

	Type := 5 //推送消息
	var routeId uint32
	routeId = crc32.ChecksumIEEE([]byte(route))

	result[0] = byte(Type)
	result[1] = byte(bodyLen >> 16)
	result[2] = byte(bodyLen >> 8)
	result[3] = byte(bodyLen >> 0)

	result[4] = byte(routeId >> 24)
	result[5] = byte(routeId >> 16)
	result[6] = byte(routeId >> 8)
	result[7] = byte(routeId >> 0)

	result = append(result, data...)

	return result
}

func BuildReplyPacket(requestId uint32, data []byte) []byte{
	if data == nil { data = make([]byte, 0)}
	var result = make([]byte, 8)
	var bodyLen = 4 + uint32(len(data))
	result[0] = byte(4)
	result[1] = byte(bodyLen >> 16)
	result[2] = byte(bodyLen >> 8)
	result[3] = byte(bodyLen >> 0)

	result[4] = byte(requestId >> 24)
	result[5] = byte(requestId >> 16)
	result[6] = byte(requestId >> 8)
	result[7] = byte(requestId >> 0)

	result = append(result, data...)

	return result
	return result
}