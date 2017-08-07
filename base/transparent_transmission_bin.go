package base

import (
	//"bytes"
	"sync"
	"sync/atomic"
)

const (
	TTB_MODE_UNAVAILABLE int32 = 0
	TTB_MODE_NORMAL      int32 = 1
)

type transparent_transmission_bin struct {
	Bin          []byte
	CheckSum     uint16
	UpgradeCount int32
	Mode         int32
}

var ttb *transparent_transmission_bin
var once sync.Once

func GetTTB() *transparent_transmission_bin {
	once.Do(func() {
		ttb = &transparent_transmission_bin{}
	})

	return ttb
}

func (tt *transparent_transmission_bin) Increase() {
	atomic.AddInt32(&tt.UpgradeCount, 1)
}

func (tt *transparent_transmission_bin) Decrease() {
	atomic.AddInt32(&tt.UpgradeCount, -1)
}

func (tt *transparent_transmission_bin) SetMode(mode int32) {
	atomic.StoreInt32(&tt.Mode, mode)
}

func (tt *transparent_transmission_bin) GetBytes(cursor int, block int) []byte {
	return tt.Bin[cursor : cursor+block]
}

func (tt *transparent_transmission_bin) GetBinSize() int {
	return len(tt.Bin)
}

func (tt *transparent_transmission_bin) IsAvailable() bool {
	return tt.Mode == TTB_MODE_NORMAL
}

func (tt *transparent_transmission_bin) IsLoadAvailable() bool {
	count := atomic.LoadInt32(&tt.UpgradeCount)
	return count == 0
}

func (tt *transparent_transmission_bin) SetBytes(bin []byte) {
	tt.Bin = bin
	tt.CheckSum = calc_check_sum(bin)
}

func calc_check_sum(bin []byte) uint16 {
	var result uint16 = 0
	for _, v := range bin {
		result += uint16(v)

	}
	//	buffer := bytes.NewBuffer(bin)
	//	buffer_count := buffer.Len()
	//	if buffer_count%2 == 1 {
	//		buffer.WriteByte(0)
	//	}
	//	var result uint16 = 0
	//	buffer_count := buffer.Len()
	//	for i := 0; i < buffer_count; i++ {
	//		result += ReadWord2(buffer)
	//	}

	return ^result
}
