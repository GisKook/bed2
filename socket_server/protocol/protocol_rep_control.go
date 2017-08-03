package protocol

import (
	"github.com/giskook/bed2/base"
)

func ParseRepControl(rep []byte) *base.SocketResponse {
	r := ParseHeader(rep)
	serialid := base.ReadDWord(r)
	result := base.ReadByte(r)

	return &base.SocketResponse{
		SerialID: serialid,
		Result:   uint32(result),
	}
}
