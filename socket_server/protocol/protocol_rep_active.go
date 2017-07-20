package protocol

import (
	"github.com/giskook/bed2/base"
)

type RepActivePkg struct {
	SerialID uint32
}

func (control *RepActivePkg) Serialize() []byte {
	return nil
}

func ParseRepActive(rep []byte) *base.SocketResponse {
	r := ParseHeader(rep)
	serialid := base.ReadDWord(r)

	return &base.SocketResponse{
		SerialID: serialid,
	}
}
