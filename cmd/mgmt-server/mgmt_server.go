package main

import (
	"github.com/chacal/thread-mgmt-server/pkg/device_registry"
	coap_routes "github.com/chacal/thread-mgmt-server/pkg/mgmt_routes/coap"
	http_routes "github.com/chacal/thread-mgmt-server/pkg/mgmt_routes/http"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/plgd-dev/go-coap/v2/mux"
	"github.com/plgd-dev/go-coap/v2/net"
	"github.com/plgd-dev/go-coap/v2/udp"
	"net/http"
	"strconv"
)

const Splash = `

----------------------------
- Thread Management Server -
----------------------------

`

type MgmtCoapServer struct {
	conn *net.UDPConn
	srv  *udp.Server
}

func NewCoapServer(coapPort int, reg *device_registry.Registry) (*MgmtCoapServer, error) {
	conn, err := net.NewListenUDP("udp6", ":"+strconv.Itoa(coapPort))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	router := mux.NewRouter()
	coap_routes.RegisterRoutes(router, reg)

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

func NewHttpServer(opts Options, reg *device_registry.Registry) *http.Server {
	router := gin.Default()
	http_routes.RegisterRoutes(router, reg)

	return &http.Server{
		Addr:    ":" + strconv.Itoa(opts.HttpPort),
		Handler: router,
	}
}
