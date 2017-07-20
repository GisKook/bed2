package protocol

import (
	"bytes"
	"github.com/giskook/bed2/base"
)

type ReqControlPkg struct {
	ID       uint64
	SerialID uint32
	Code     uint8
	Action   uint8
}

func (p *ReqControlPkg) Serialize() []byte {
	var writer bytes.Buffer
	WriteHeader(&writer, PROTOCOL_REP_CONTROL)
	base.WriteDWord(&writer, p.SerialID)
	base.WriteByte(&writer, p.Code)
	base.WriteByte(&writer, p.Action)
	WriteTail(&writer)

	return writer.Bytes()
}

func ParseReqControl(req *base.HttpRequest) *ReqControlPkg {
	return &ReqControlPkg{
		ID:       req.BedID,
		SerialID: req.SerialID,
		Code:     req.Code,
		Action:   req.Action,
	}
}
