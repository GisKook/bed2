package protocol

import (
	"github.com/giskook/bed2/base"
)

func ParseRepActive(rep []byte) *base.SocketResponse {
	r := ParseHeader(rep)
	serialid := base.ReadDWord(r)

	return &base.SocketResponse{
		SerialID: serialid,
	}
}
