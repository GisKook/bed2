package http_server

import (
	"encoding/json"
	"fmt"
	"github.com/giskook/bed2/base"
	"net/http"
	"time"
)

type TTCode int

const (
	TT_ESTABLISHED     TTCode = 0
	TT_UNAVAILABLE     TTCode = 1
	TT_LACK_PARAMTER   TTCode = 2
	TT_INTERAL_ERROR   TTCode = 3
	TT_TIMEOUT         TTCode = 4
	TT_BIN_UNAVAILABLE TTCode = 5
)

var CODE_DESC = []string{"成功, 正在升级中，请稍等", "不可用", "请检查参数", "服务器内部错误", "超时", "透传文件不可用"}

func (c TTCode) Desc() string {
	return CODE_DESC[int(c)]
}

type RepTransparentTransmission struct {
	Code int    `json:"code"`
	Desc string `json:"desc"`
}

func EncodeTransparentTransmission(code TTCode, r *http.Request) string {
	resp := &RepTransparentTransmission{
		Code: int(code),
		Desc: code.Desc(),
	}

	response, _ := json.Marshal(resp)
	RecordRecv(r, string(response))

	return string(response)
}

func (h *HttpServer) HandleTransparentTransmission(endpoint string) {
	h.mux.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if x := recover(); x != nil {
				fmt.Fprint(w, EncodeTransparentTransmission(TT_INTERAL_ERROR, r))
			}
		}()
		RecordSend(r)
		r.ParseForm()
		mac := r.Form.Get("mac")
		if mac == "" || len(mac) != 12 {
			fmt.Fprint(w, EncodeTransparentTransmission(TT_LACK_PARAMTER, r))
			RecordSendLackParamter(r)
			return
		}
		if !base.GetTTB().IsAvailable() {
			fmt.Fprint(w, EncodeTransparentTransmission(TT_BIN_UNAVAILABLE, r))
			return
		}
		_mac, err := base.GetMac(mac)
		if err == nil {
			ch, index := h.makeRequest()
			h.request(&base.HttpRequest{
				BedID:     _mac,
				SerialID:  index,
				RequestID: base.PROTOCOL_REQ_TRANSPARNET_TRANSMISSION,
				ByteCount: uint32(base.GetTTB().GetBinSize()),
			})

			select {
			case resp := <-ch:
				if resp.Result == 0 {
					fmt.Fprint(w, EncodeTransparentTransmission(TT_ESTABLISHED, r))
				} else {
					fmt.Fprint(w, EncodeTransparentTransmission(TT_UNAVAILABLE, r))
				}

			case <-time.After(time.Duration(h.conf.TimeOut) * time.Second):
				fmt.Fprint(w, EncodeTransparentTransmission(TT_TIMEOUT, r))
			}
			h.router.delRequest(index)
		}
	})
}
