package protocol

import (
	"bytes"
	"github.com/giskook/bed2/base"
)

type ReqLoginPkg struct {
	ID              uint64
	HardwareVersion uint8
	SoftwareVersion uint8
}

func (login *ReqLoginPkg) Serialize() []byte {
	var writer bytes.Buffer
	WriteHeader(&writer, PROTOCOL_REP_LOGIN)
	base.WriteByte(&writer, 0)

	WriteTail(&writer)

	return writer.Bytes()
}

func ParseReqLogin(buf []byte) *ReqLoginPkg {
	reader := ParseHeader(buf)

	return &ReqLoginPkg{
		ID:              base.ReadMac(reader),
		HardwareVersion: base.ReadByte(reader),
		SoftwareVersion: base.ReadByte(reader),
	}
}
