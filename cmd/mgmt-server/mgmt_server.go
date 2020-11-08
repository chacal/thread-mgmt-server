package main

import (
	"github.com/chacal/thread-mgmt-server/pkg/device_registry"
	"github.com/chacal/thread-mgmt-server/pkg/mgmt_routes"
	"github.com/pkg/errors"
	"github.com/plgd-dev/go-coap/v2/mux"
	"github.com/plgd-dev/go-coap/v2/net"
	"github.com/plgd-dev/go-coap/v2/udp"
	"strconv"
)

type MgmtCoapServer struct {
	conn *net.UDPConn
	srv  *udp.Server
}

func NewServer(opts Options, reg *device_registry.Registry) (*MgmtCoapServer, error) {
	conn, err := net.NewListenUDP("udp6", ":"+strconv.Itoa(opts.Port))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	router := mux.NewRouter()
	mgmt_routes.RegisterRoutes(router, reg)

	srv := udp.NewServer(udp.WithMux(router), udp.WithKeepAlive(nil))

	return &MgmtCoapServer{conn, srv}, nil
}

func (s *MgmtCoapServer) Serve() error {
	return s.srv.Serve(s.conn)
}

func (s *MgmtCoapServer) Stop() {
	s.srv.Stop()
	s.conn.Close()
}
