package protocol

import (
	"bytes"
	"github.com/giskook/bed2/base"
)

type ReqHeartPkg struct {
	ID uint64
}

func (h *ReqHeartPkg) Serialize() []byte {
	var writer bytes.Buffer
	WriteHeader(&writer, PROTOCOL_REP_HEART)
	WriteTail(&writer)

	return writer.Bytes()
}

func ParseReqHeart(buf []byte) *ReqHeartPkg {
	reader := ParseHeader(buf)

	return &ReqHeartPkg{
		ID: base.ReadMac(reader),
	}
}
