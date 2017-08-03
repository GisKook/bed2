package http_server

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
)

type Code int

const (
	HTTP_CONSTANT_SEND_REQ          string = "<http req> %x %s\n"
	HTTP_CONSTANT_SEND_REP          string = "<http rep> %x %s\n"
	HTTP_CONSTANT_SEND_REQ_LACK_PAR string = "请检查参数\n"

	HTTP_REP_SUCCESS        Code = 0
	HTTP_REP_LACK_PARAMETER Code = 1
	HTTP_REP_INTERAL_ERROR  Code = 2
	HTTP_REP_TIMEOUT        Code = 3
)

var HTTP_REQUEST_DESC []string = []string{
	"成功",
	"缺少参数",
	"服务器内部错误",
	"超时"}

func (c Code) Desc() string {
	return HTTP_REQUEST_DESC[c]
}

type GeneralResponse struct {
	Code int    `json:"code"`
	Desc string `json:"desc"`
}

func EncodeGeneralResponse(code Code) string {
	general_response := &GeneralResponse{
		Code: int(code),
		Desc: code.Desc(),
	}

	resp, _ := json.Marshal(general_response)

	return string(resp)
}

func RecordSend(r *http.Request) {
	v, e := httputil.DumpRequest(r, true)
	if e != nil {
		log.Println(e.Error())
		return
	}
	log.Println(string(v))
}

func RecordRecv(r *http.Request, resp string) {
	log.Printf(HTTP_CONSTANT_SEND_REP, &r.URL, resp)
}

func RecordSendLackParamter(r *http.Request) {
	log.Println(HTTP_CONSTANT_SEND_REQ_LACK_PAR)
}
