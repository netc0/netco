package connector

import (
	"time"
)

func (s *Session) procPacket(packet Packet) int {
	if packet.Type == PacketType_SYN { // 收到 SYN
		s.procSYN(packet)
		return 0
	} else if packet.Type == PacketType_ACK { // 收到 ACK
		s.procACK(packet)
		return 0
	} else if packet.Type == PacketType_HEARTBEAT { // 纯心跳包 一般不需要
		s.procHeartBeat(packet)
		return 0
	} else if packet.Type == PacketType_DATA { // on data
		s.procData(packet)
		return 0
	} else if packet.Type == PacketType_KICK { // on kick

	}

	return -1
}

// 处理 SYN
func (s *Session) procSYN(packet Packet) {
	s.Push(PacketType_ACK, nil) // 必须回应SYN
}

// 处理 ACK
func (this *Session) procACK(packet Packet) { // on connected
	//this.transporter.onConnectionData(this, packet.Body)
}

// 处理 心跳包
func (s *Session) procHeartBeat(packet Packet) {
	s.heartBeatTime = time.Now()
}

// 处理数据包
func (this *Session) procData(packet Packet) {
	// [requestId] [routeId] [data]
	// 1. 解析出requestId
	var requestId uint32
	var routeId uint32
	var data = packet.Body

	requestId = uint32(
		uint32(data[0]) << 24 |
		uint32(data[1]) << 16 |
		uint32(data[2]) << 8 |
		uint32(data[3]));
	// 2. 解析出 routeId
	routeId = uint32(
		uint32(data[4]) << 24 |
		uint32(data[5]) << 16 |
		uint32(data[6]) << 8 |
		uint32(data[7]));
	data = data[8:]
	this.transporter.onConnectionData(this, requestId, routeId, data)
}

// 处理 Kick
func (s *Session) procKick(packet Packet) {
	s.conn.Close()
}