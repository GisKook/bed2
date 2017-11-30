package protocol

import (
	"bytes"
	"encoding/binary"
	"github.com/giskook/bed2/base"
)

const (
	PROTOCOL_START_FLAG    byte   = 0xce
	PROTOCOL_END_FLAG      byte   = 0xce
	PROTOCOL_BED_RECV_FLAG uint16 = 0x0bed
	PROTOCOL_BED_SEND_FLAG uint16 = 0x8bed

	PROTOCOL_ILLEGAL   uint16 = 0xffff
	PROTOCOL_HALF_PACK uint16 = 0xfffe

	PROTOCOL_MIN_LEN uint16 = 9
	PROTOCOL_MAX_LEN uint16 = 17

	PROTOCOL_REQ_LOGIN uint16 = 0x8001
	PROTOCOL_REP_LOGIN uint16 = 0x0001
	PROTOCOL_REQ_HEART uint16 = 0x8002
	PROTOCOL_REP_HEART uint16 = 0x0002

	PROTOCOL_REQ_CONTROL     uint16 = 0x0003
	PROTOCOL_REP_CONTROL     uint16 = 0x8003
	PROTOCOL_REQ_ACTIVE_TEST uint16 = 0x0004
	PROTOCOL_REP_ACTIVE_TEST uint16 = 0x8004

	PROTOCOL_REQ_TRANSPARENT_TRANSMISSION = 0x0005
	PROTOCOL_REP_TRANSPARENT_TRANSMISSION = 0x8005

	PROTOCOL_REQ_BIN_FILE = 0x8006

	PROTOCOL_REP_TRANSPARENT_TRANSMISSION_STOP = 0x8007
)

func ParseHeader(buffer []byte) *bytes.Reader {
	reader := bytes.NewReader(buffer)
	reader.Seek(7, 0)

	return reader
}

func WriteHeader(writer *bytes.Buffer, cmdid uint16) {
	writer.WriteByte(PROTOCOL_START_FLAG)
	base.WriteWord(writer, 0)
	base.WriteWord(writer, PROTOCOL_BED_RECV_FLAG)
	base.WriteWord(writer, cmdid)
}

func WriteTail(writer *bytes.Buffer) {
	WriteLength(writer)
	sum := checkSum(writer.Bytes(), uint16(writer.Len()))
	base.WriteByte(writer, sum)
	base.WriteByte(writer, PROTOCOL_END_FLAG)
}

func WriteLength(writer *bytes.Buffer) {
	length := writer.Len()
	length += 2
	length_byte := make([]byte, 2)
	binary.BigEndian.PutUint16(length_byte, uint16(length))
	writer.Bytes()[1] = length_byte[0]
	writer.Bytes()[2] = length_byte[1]
}

func checkSum(cmd []byte, cmdlen uint16) byte {
	temp := cmd[0]
	for i := uint16(1); i < cmdlen; i++ {
		temp ^= cmd[i]
	}

	return temp
}

func CheckProtocol(buffer *bytes.Buffer) (uint16, uint16) {
	bufferlen := buffer.Len()
	if bufferlen == 0 {
		return PROTOCOL_ILLEGAL, 0
	}
	if buffer.Bytes()[0] != PROTOCOL_START_FLAG {
		buffer.ReadByte()
		CheckProtocol(buffer)
	} else if bufferlen > 2 {
		len_in_protocol := base.GetWord(buffer.Bytes()[1:3])
		if len_in_protocol < PROTOCOL_MIN_LEN || len_in_protocol > PROTOCOL_MAX_LEN {
			buffer.ReadByte()
			CheckProtocol(buffer)
		}

		if int(len_in_protocol) > bufferlen {
			return PROTOCOL_HALF_PACK, 0
		} else {
			//chkSum := checkSum(buffer.Bytes(), len_in_protocol-2)
			//if chkSum == buffer.Bytes()[len_in_protocol-2] && buffer.Bytes()[len_in_protocol-1] == PROTOCOL_END_FLAG {
			if buffer.Bytes()[len_in_protocol-1] == PROTOCOL_END_FLAG {
				protocol_id := base.GetWord(buffer.Bytes()[5:7])
				return protocol_id, len_in_protocol
			} else {
				buffer.ReadByte()
				CheckProtocol(buffer)
			}
		}
	} else {
		return PROTOCOL_HALF_PACK, 0
	}

	return PROTOCOL_HALF_PACK, 0
}
