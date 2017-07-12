package protocol

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
)

//func ParseHeader(buffer []byte) (*bytes.Reader, uint16, uint16, uint64) {
//	reader := bytes.NewReader(buffer)
//	reader.Seek(1, 0)
//	length := base.ReadWord(reader)
//	protocol_id := base.ReadWord(reader)
//	tid_string := base.ReadBcdString(reader, 8)
//	tid, _ := strconv.ParseUint(tid_string, 10, 64)
//
//	return reader, length, protocol_id, tid
//}
//
//func WriteHeader(writer *bytes.Buffer, length uint16, cmdid uint16, cpid uint64) {
//	writer.WriteByte(PROTOCOL_START_FLAG)
//	base.WriteWord(writer, length)
//	base.WriteWord(writer, cmdid)
//	base.WriteBcdCpid(writer, cpid)
//}
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
			chkSum := checkSum(buffer.Bytes(), len_in_protocol-2)
			if chkSum == buffer.Bytes()[len_in_protocol-2] && buffer.Bytes()[len_in_protocol-1] == PROTOCOL_END_FLAG {
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
