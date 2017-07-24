package protocol

import (
	"bytes"
	"github.com/giskook/bed2/base"
)

type ReqActivePkg struct {
	ID       uint64
	SerialID uint32
}

func (control *ReqActivePkg) Serialize() []byte {
	var writer bytes.Buffer
	WriteHeader(&writer, PROTOCOL_REQ_ACTIVE_TEST)
	base.WriteDWord(&writer, control.SerialID)
	WriteTail(&writer)

	return writer.Bytes()
}

func ParseReqActiveTest(req *base.HttpRequest) *ReqActivePkg {
	return &ReqActivePkg{
		ID:       req.BedID,
		SerialID: req.SerialID,
	}
}
