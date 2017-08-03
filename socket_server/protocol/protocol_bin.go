package protocol

import (
	"github.com/giskook/bed2/base"
)

type ReqBinPkg struct {
	Cursor    uint16
	BlockSize uint16
}

func ParseReqBin(buf []byte) *ReqBinPkg {
	reader := ParseHeader(buf)

	return &ReqBinPkg{
		Cursor:    base.ReadWord(reader),
		BlockSize: base.ReadWord(reader),
	}
}
