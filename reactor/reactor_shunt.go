package reactor

import (
	"github.com/giskook/bed2/base"
	"github.com/giskook/bed2/socket_server/protocol"
)

func (r *Reactor) send_offline_pkg(index uint32) {
	r.http.Response(&base.SocketResponse{
		SerialID: index,
		Result:   base.PROTOCOL_BED_OFFLINE,
	})
}

func (r *Reactor) shunt() {
	for {
		select {
		case <-r.exit:
			return
		case req := <-r.http.HttpRequestChan:
			if req.RequestID == base.PROTOCOL_REQ_CONTROL {
				p := protocol.ParseReqControl(req)
				err := r.ss.Send(req.BedID, p)
				if err == base.ErrBedOffline {
					r.send_offline_pkg(p.SerialID)
				}
			} else if req.RequestID == base.PROTOCOL_REQ_ACTIVE_TEST {
				p := protocol.ParseReqActiveTest(req)
				err := r.ss.Send(req.BedID, p)
				if err == base.ErrBedOffline {
					r.send_offline_pkg(p.SerialID)
				}
			} else if req.RequestID == base.PROTOCOL_REQ_TRANSPARNET_TRANSMISSION {
				p := protocol.ParseReqTransparentTransmission(req)
				err := r.ss.Send(req.BedID, p)
				if err == base.ErrBedOffline {
					r.send_offline_pkg(p.SerialID)
				}
			}
		case rep := <-r.ss.SocketOutResult:
			r.http.Response(rep)
		}
	}
}
