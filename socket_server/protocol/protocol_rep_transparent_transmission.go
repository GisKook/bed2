package protocol

import (
	"github.com/giskook/bed2/base"
)

func ParseRepTransparentTransmission(rep []byte) *base.SocketResponse {
	r := ParseHeader(rep)
	return &base.SocketResponse{
		SerialID: base.ReadDWord(r),
		Result:   uint32(base.ReadByte(r)),
	}
}
