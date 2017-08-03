package base

const (
	PROTOCOL_REQ_CONTROL                  uint16 = 0x0003
	PROTOCOL_REQ_ACTIVE_TEST              uint16 = 0x0004
	PROTOCOL_REQ_TRANSPARNET_TRANSMISSION uint16 = 0x0005

	PROTOCOL_BED_OFFLINE uint32 = 4

	CTL_CODE_UNSUPPORT      uint8 = 0x00
	CTL_CODE_BEDPAN_OPEN    uint8 = 0x01 //便盆开
	CTL_CODE_BEDPAN_CLOSE   uint8 = 0x02 //便盆关
	CTL_CODE_TURN_LEFT      uint8 = 0x03 //左侧翻
	CTL_CODE_TURN_RIGHT     uint8 = 0x04 //右侧翻
	CTL_CODE_RAISE_BACK     uint8 = 0x05 //起背
	CTL_CODE_LAND_BACK      uint8 = 0x06 //背平
	CTL_CODE_RAISE_LEG      uint8 = 0x07 //抬腿
	CTL_CODE_LAND_LEG       uint8 = 0x08 //落腿
	CTL_CODE_AUTO_TURN      uint8 = 0x09 //自动翻
	CTL_CODE_SAT_UP         uint8 = 0x0a //坐起
	CTL_CODE_EMERGENCY_STOP uint8 = 0x0b //急停
	CTL_CODE_RESET          uint8 = 0x0c //复位

	CTL_ACTION_MOVE uint8 = 1
	CTL_ACTION_STOP uint8 = 0
)

type HttpRequest struct {
	BedID     uint64
	SerialID  uint32
	RequestID uint16
	Code      uint8
	Action    uint8
	Result    uint32
	ByteCount uint32
	CheckSum  uint16
}

type SocketResponse struct {
	BedID    uint64
	SerialID uint32
	Result   uint32
}
