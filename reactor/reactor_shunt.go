package reactor

import (
	"github.com/giskook/bed2/base"
	"github.com/giskook/bed2/socket_server/protocol"
)

func (r *Reactor) shunt() {
	for {
		select {
		case <-r.exit:
		case req := <-r.http.HttpRequestChan:
			if req.RequestID == base.PROTOCOL_REQ_CONTROL {
				p := protocol.ParseReqControl(req)
				r.ss.Send(req.BedID, p)
			} else if req.RequestID == base.PROTOCOL_REQ_ACTIVE_TEST {
				p := protocol.ParseReqActiveTest(req)
				r.ss.Send(req.BedID, p)
			}
		case rep := <-r.ss.SocketOutResult:
			r.http.Response(rep)
		}
	}
}
