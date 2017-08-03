package protocol

import (
	"bytes"
	"github.com/giskook/bed2/base"
)

type ReqTransparentTransmissionPkg struct {
	ID        uint64
	SerialID  uint32
	ByteCount uint32
	CheckSum  uint16
}

func (p *ReqTransparentTransmissionPkg) Serialize() []byte {
	var writer bytes.Buffer
	WriteHeader(&writer, PROTOCOL_REQ_TRANSPARENT_TRANSMISSION)
	base.WriteDWord(&writer, p.SerialID)
	base.WriteDWord(&writer, p.ByteCount)
	base.WriteWord(&writer, p.CheckSum)
	WriteTail(&writer)

	return writer.Bytes()
}

func ParseReqTransparentTransmission(req *base.HttpRequest) *ReqTransparentTransmissionPkg {
	return &ReqTransparentTransmissionPkg{
		ID:        req.BedID,
		SerialID:  req.SerialID,
		ByteCount: req.ByteCount,
		CheckSum:  req.CheckSum,
	}
}
