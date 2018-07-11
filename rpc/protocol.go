package rpc

import (
	"hash/crc32"
)

type mailProtocol struct {
	readType int
	buffer []byte

	mailType uint32
	mailLength uint32
}

func (this*mailProtocol) Init() {
	this.readType = 0
}

func (this *mailProtocol) handleBytes(data []byte) *Mail {
	if data != nil {
		this.buffer = append(this.buffer, data...)
	}
	if this.readType == 0 { // header
		if len(this.buffer) < 8 {
			return nil
		}
		// 解析header
		this.mailType = uint32(
			uint32(this.buffer[0]) << 24 |
			uint32(this.buffer[1]) << 16 |
			uint32(this.buffer[2]) << 8 |
			uint32(this.buffer[3]))
		this.mailLength = uint32(
			uint32(this.buffer[4]) << 24 |
			uint32(this.buffer[5]) << 16 |
			uint32(this.buffer[6]) << 8 |
			uint32(this.buffer[7]))
		this.buffer = this.buffer[8:]
		this.readType = 1
		return this.handleBytes(nil)
	} else {
		if uint32(len(this.buffer)) >= this.mailLength { // 解析 body 完毕
			var np Mail
			np.Type = this.mailType
			if this.mailLength > 0 {
				np.data = this.buffer[:this.mailLength]
				this.buffer = this.buffer[this.mailLength:]
			}

			this.mailLength = 0
			this.readType = 0
			return &np
		}
	}
	return nil
}

func (this* mailProtocol) encode(mail *Mail) []byte {
	var result []byte
	result = make([]byte, 8)

	result[0] = byte(mail.Type >> 24)
	result[1] = byte(mail.Type >> 16)
	result[2] = byte(mail.Type >> 8)
	result[3] = byte(mail.Type >> 0)

	mail.Encode()
	size := len(mail.data)
	result[4] = byte(size >> 24)
	result[5] = byte(size >> 16)
	result[6] = byte(size >> 8)
	result[7] = byte(size >> 0)

	result = append(result, mail.data...)

	return result
}

func RouteHash(route string) uint32 {
	return crc32.ChecksumIEEE([]byte(route))
}