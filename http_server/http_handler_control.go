package http_server

import (
	"encoding/json"
	"fmt"
	"github.com/giskook/bed2/base"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	HTTP_CTL_FAIL          int = 0
	HTTP_CTL_SUCCESS       int = 1
	HTTP_CTL_LACK_PARAMTER int = 2
	HTTP_CTL_INTERAL_ERROR int = 3

	CTL_BEDPAN_OPEN_MOVE  string = "10"
	CTL_BEDPAN_OPEN_STOP  string = "11"
	CTL_BEDPAN_CLOSE_MOVE string = "20"
	CTL_BEDPAN_CLOSE_STOP string = "21"
	CTL_TURN_LEFT_MOVE    string = "30"
	CTL_TURN_LEFT_STOP    string = "31"
	CTL_TURN_RIGHT_MOVE   string = "40"
	CTL_TURN_RIGHT_STOP   string = "41"
	CTL_RAISE_BACK_MOVE   string = "50"
	CTL_RAISE_BACK_STOP   string = "51"
	CTL_LAND_BACK_MOVE    string = "60"
	CTL_LAND_BACK_STOP    string = "61"
	CTL_RAISE_LEG_MOVE    string = "70"
	CTL_RAISE_LEG_STOP    string = "71"
	CTL_LAND_LEG_MOVE     string = "80"
	CTL_LAND_LEG_STOP     string = "81"
	CTL_AUTO_TURN         string = "9"
	CTL_SAT_UP            string = "0a"
	CTL_EMERGENCY_STOP    string = "0b"
	CTL_RESET             string = "0c"
)

var CTL_CODE []string = []string{
	CTL_BEDPAN_OPEN_MOVE,
	CTL_BEDPAN_OPEN_STOP,
	CTL_BEDPAN_CLOSE_MOVE,
	CTL_BEDPAN_CLOSE_STOP,
	CTL_TURN_LEFT_MOVE,
	CTL_TURN_LEFT_STOP,
	CTL_TURN_RIGHT_MOVE,
	CTL_TURN_RIGHT_STOP,
	CTL_RAISE_BACK_MOVE,
	CTL_RAISE_BACK_STOP,
	CTL_LAND_BACK_MOVE,
	CTL_LAND_BACK_STOP,
	CTL_RAISE_LEG_MOVE,
	CTL_RAISE_LEG_STOP,
	CTL_LAND_LEG_MOVE,
	CTL_LAND_LEG_STOP,
	CTL_AUTO_TURN,
	CTL_SAT_UP,
	CTL_EMERGENCY_STOP,
	CTL_RESET}

func isValidParamter(code string) bool {
	for _, v := range CTL_CODE {
		if code == v {
			return true
		}
	}

	return false
}

func ParseCode(code string) (uint8, uint8) {
	switch code {
	case CTL_BEDPAN_OPEN_MOVE:
		return base.CTL_CODE_BEDPAN_OPEN, base.CTL_ACTION_MOVE
	case CTL_BEDPAN_OPEN_STOP:
		return base.CTL_CODE_BEDPAN_OPEN, base.CTL_ACTION_STOP
	case CTL_BEDPAN_CLOSE_MOVE:
		return base.CTL_CODE_BEDPAN_CLOSE, base.CTL_ACTION_MOVE
	case CTL_BEDPAN_CLOSE_STOP:
		return base.CTL_CODE_BEDPAN_CLOSE, base.CTL_ACTION_STOP
	case CTL_TURN_LEFT_MOVE:
		return base.CTL_CODE_TURN_LEFT, base.CTL_ACTION_MOVE
	case CTL_TURN_LEFT_STOP:
		return base.CTL_CODE_TURN_LEFT, base.CTL_ACTION_STOP
	case CTL_TURN_RIGHT_MOVE:
		return base.CTL_CODE_TURN_RIGHT, base.CTL_ACTION_MOVE
	case CTL_TURN_RIGHT_STOP:
		return base.CTL_CODE_TURN_RIGHT, base.CTL_ACTION_STOP
	case CTL_RAISE_BACK_MOVE:
		return base.CTL_CODE_RAISE_BACK, base.CTL_ACTION_MOVE
	case CTL_RAISE_BACK_STOP:
		return base.CTL_CODE_RAISE_BACK, base.CTL_ACTION_STOP
	case CTL_LAND_BACK_MOVE:
		return base.CTL_CODE_LAND_BACK, base.CTL_ACTION_MOVE
	case CTL_LAND_BACK_STOP:
		return base.CTL_CODE_LAND_BACK, base.CTL_ACTION_STOP
	case CTL_RAISE_LEG_MOVE:
		return base.CTL_CODE_RAISE_LEG, base.CTL_ACTION_MOVE
	case CTL_RAISE_LEG_STOP:
		return base.CTL_CODE_RAISE_LEG, base.CTL_ACTION_STOP
	case CTL_LAND_LEG_MOVE:
		return base.CTL_CODE_LAND_LEG, base.CTL_ACTION_MOVE
	case CTL_LAND_LEG_STOP:
		return base.CTL_CODE_LAND_LEG, base.CTL_ACTION_STOP
	case CTL_AUTO_TURN:
		return base.CTL_CODE_AUTO_TURN, base.CTL_ACTION_MOVE
	case CTL_SAT_UP:
		return base.CTL_CODE_SAT_UP, base.CTL_ACTION_MOVE
	case CTL_EMERGENCY_STOP:
		return base.CTL_CODE_EMERGENCY_STOP, base.CTL_ACTION_MOVE
	case CTL_RESET:
		return base.CTL_CODE_RESET, base.CTL_ACTION_MOVE
	}

	return base.CTL_CODE_UNSUPPORT, base.CTL_ACTION_MOVE
}

type RepControl struct {
	Result int `json:"result"`
}

func EncodeControl(code int, r *http.Request) string {
	resp := &RepControl{
		Result: code,
	}

	response, _ := json.Marshal(resp)
	RecordRecv(r, string(response))

	return string(response)
}

func (h *HttpServer) HandleControl(endpoint string) {
	h.mux.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if x := recover(); x != nil {
				fmt.Fprint(w, EncodeControl(HTTP_CTL_INTERAL_ERROR, r))
			}
		}()
		RecordSend(r)
		r.ParseForm()
		mac := r.Form.Get("mac")
		_code := r.Form.Get("type")
		code := strings.ToLower(_code)
		if mac == "" || len(mac) != 12 || !isValidParamter(code) {
			fmt.Fprint(w, EncodeControl(HTTP_CTL_LACK_PARAMTER, r))
			RecordSendLackParamter(r)
			return
		}
		_mac, err := base.GetMac(mac)
		if err == nil {
			ch, index := h.makeRequest()
			_code, action := ParseCode(code)

			h.request(&base.HttpRequest{
				BedID:     _mac,
				SerialID:  index,
				RequestID: base.PROTOCOL_REQ_CONTROL,
				Code:      _code,
				Action:    action,
			})

			select {
			case resp := <-ch:
				fmt.Fprint(w, EncodeControl(int(resp.Result), r))

			case <-time.After(time.Duration(h.conf.TimeOut) * time.Second):
				fmt.Fprint(w, EncodeControl(HTTP_CTL_INTERAL_ERROR, r))
			}
			h.router.delRequest(index)
		} else {
			log.Println(err.Error())
		}
	})
}
