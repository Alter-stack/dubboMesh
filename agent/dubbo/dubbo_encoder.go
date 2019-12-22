package dubbo

import (
	"dubboMesh/agent/dubbo/server"
	"encoding/binary"
)

/**
 * 协议头是16字节的定长数据
 * 2字节magic字符串0xdabb,0-7高位，8-15低位
 * 1字节的消息标志位。16-20序列id,21 event,22 two way,23请求或响应标识
 * 1字节状态。当消息类型为响应时，设置响应状态。24-31位。
 * 8字节，消息ID,long类型，32-95位。
 * 4字节，消息长度，96-127位
 **/
const (
	// header length.
	HEADER_LENGTH = 16

	// magic header
	MAGIC = uint16(0xdabb)
	MAGIC_HIGH = byte(0xda)
	MAGIC_LOW = byte(0xbb)

	// message flag
	FLAG_REQUEST = byte(0x80) // 10000000
	FLAG_TWOWAY  = byte(0x40) // 01000000
	FLAG_EVENT   = byte(0x20) // for heartbeat

	SERIALIZATION_MASK = 0x1f

	DUBBO_VERSION = "2.0.1"
)

func PackRequest(req *server.AgentRequest) ([]byte, error) {
	Buf := make([]byte, 2048)
	totalBuf := Buf[0:0]
	startIndex := len(totalBuf)
	// FLAG_REQUEST | FLAG_TWOWAY = 11000000
	dubboHeader := [HEADER_LENGTH]byte{MAGIC_HIGH, MAGIC_LOW, FLAG_REQUEST | FLAG_TWOWAY}
	// magic
	totalBuf = append(totalBuf, dubboHeader[:]...)

	// status, 仅在 Req/Res 为0（响应）时有用，用于标识响应的状态
	totalBuf[3+startIndex] = 0

	// Serialization ID
	// serialization id, two way flag, event, request/response flag, 比如 fastjson 的值为6
	totalBuf[2+startIndex] |= byte(FLAG_REQUEST | 6)

	// Request ID
	binary.LittleEndian.PutUint64(totalBuf[4+startIndex:], uint64(req.RequestID))

	data := encodeRequestData(req)
	mdataLen := len(data)
	binary.LittleEndian.PutUint32(totalBuf[4+8+startIndex:], uint32(mdataLen))

	nowLen := len(totalBuf)
	buf := Buf[nowLen:]
	copy(buf, data[:])
	totalBuf = append(totalBuf, buf[:]...)
	return totalBuf, nil
}

func encodeRequestData(req *server.AgentRequest) []byte {
	return []byte{1,2}
}