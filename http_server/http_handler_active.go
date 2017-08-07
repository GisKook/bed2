package http_server

import (
	"encoding/json"
	"fmt"
	"github.com/giskook/bed2/base"
	"net/http"
	"time"
)

const (
	HTTP_NOT_CONNECTED int = 0
	HTTP_CONNECTED     int = 1
	HTTP_LACK_PARAMTER int = 2
	HTTP_INTERAL_ERROR int = 3
)

type RepActiveTest struct {
	CheckOnline int `json:"checkOnline"`
}

func EncodeActiveTest(code int, r *http.Request) string {
	resp := &RepActiveTest{
		CheckOnline: code,
	}

	response, _ := json.Marshal(resp)
	RecordRecv(r, string(response))

	return string(response)
}

func (h *HttpServer) ActiveTestHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if x := recover(); x != nil {
			fmt.Fprint(w, EncodeActiveTest(HTTP_INTERAL_ERROR, r))
		}
	}()
	RecordSend(r)
	r.ParseForm()
	mac := r.Form.Get("mac")
	if mac == "" || len(mac) != 12 {
		fmt.Fprint(w, EncodeActiveTest(HTTP_LACK_PARAMTER, r))
		RecordSendLackParamter(r)
		return
	}
	_mac, err := base.GetMac(mac)
	if err == nil {
		ch, index := h.makeRequest()
		h.request(&base.HttpRequest{
			BedID:     _mac,
			SerialID:  index,
			RequestID: base.PROTOCOL_REQ_ACTIVE_TEST,
		})

		select {
		case resp := <-ch:
			if resp.Result == 0 {
				fmt.Fprint(w, EncodeActiveTest(HTTP_CONNECTED, r))
			} else {
				fmt.Fprint(w, EncodeActiveTest(HTTP_NOT_CONNECTED, r))
			}

		case <-time.After(time.Duration(h.conf.TimeOut) * time.Second):
			fmt.Fprint(w, EncodeActiveTest(HTTP_NOT_CONNECTED, r))
		}
		h.router.delRequest(index)
	}
}

func (h *HttpServer) HandleActiveTest(endpoint string) {
	h.mux.HandleFunc(endpoint, h.ActiveTestHandler)
	//	h.mux.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
	//		defer func() {
	//			if x := recover(); x != nil {
	//				fmt.Fprint(w, EncodeActiveTest(HTTP_INTERAL_ERROR, r))
	//			}
	//		}()
	//		RecordSend(r)
	//		r.ParseForm()
	//		mac := r.Form.Get("mac")
	//		if mac == "" || len(mac) != 12 {
	//			fmt.Fprint(w, EncodeActiveTest(HTTP_LACK_PARAMTER, r))
	//			RecordSendLackParamter(r)
	//			return
	//		}
	//		_mac, err := base.GetMac(mac)
	//		if err == nil {
	//			ch, index := h.makeRequest()
	//			h.request(&base.HttpRequest{
	//				BedID:     _mac,
	//				SerialID:  index,
	//				RequestID: base.PROTOCOL_REQ_ACTIVE_TEST,
	//			})
	//
	//			select {
	//			case resp := <-ch:
	//				if resp.Result == 0 {
	//					fmt.Fprint(w, EncodeActiveTest(HTTP_CONNECTED, r))
	//				} else {
	//					fmt.Fprint(w, EncodeActiveTest(HTTP_NOT_CONNECTED, r))
	//				}
	//
	//			case <-time.After(time.Duration(h.conf.TimeOut) * time.Second):
	//				fmt.Fprint(w, EncodeActiveTest(HTTP_NOT_CONNECTED, r))
	//			}
	//			h.router.delRequest(index)
	//		}
	//	})
}
