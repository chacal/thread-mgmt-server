package main

import (
	"github.com/chacal/thread-mgmt-server/pkg/coap_utils"
	"github.com/chacal/thread-mgmt-server/pkg/server"
	"github.com/pkg/errors"
	"github.com/plgd-dev/go-coap/v2/mux"
	"github.com/plgd-dev/go-coap/v2/net"
	"github.com/plgd-dev/go-coap/v2/udp"
	log "github.com/sirupsen/logrus"
	gonet "net"
	"strconv"
)

type Options struct {
	Interface         string `short:"i" long:"interface" description:"Interface to bind to" default:"wpan0" env:"BIND_INTERFACE"`
	ListenAddr        string `short:"l" long:"listen" description:"Multicast address to listen" default:"[ff03::1]" env:"LISTEN_ADDRESS"`
	Port              int    `short:"p" long:"port" description:"Port to listen" default:"5683" env:"PORT"`
	MgmtServerAddress string `short:"m" long:"mgmt-server" description:"Address of the returned management server" env:"MGMT_SERVER_ADDRESS" required:"yes"`
}

type DiscoveryResponse struct {
	MgmtServer string `json:"mgmtServer,omitempty"`
}

func main() {
	log.SetFormatter(&log.TextFormatter{ForceColors: true, DisableTimestamp: true})
	opts := Options{}
	server.ParseOptions(&opts)
	LogOptions(opts)
	log.Info(Splash)
	err := startCoapServer(opts)
	log.Fatalf("%+v", err)
}

func startCoapServer(opts Options) error {
	conn, err := createServerConn(opts)
	if err != nil {
		return errors.WithStack(err)
	}
	defer conn.Close()

	router := mux.NewRouter()
	router.Use(coap_utils.LoggingMiddleware)
	router.Handle("/discover", handleGetDiscover(opts.MgmtServerAddress))

	server := udp.NewServer(udp.WithMux(router), udp.WithKeepAlive(nil))
	defer server.Stop()
	return server.Serve(conn)
}

func createServerConn(opts Options) (*net.UDPConn, error) {
	listenAddrWithPort := opts.ListenAddr + ":" + strconv.Itoa(opts.Port)

	conn, err := net.NewListenUDP("udp6", listenAddrWithPort)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	err = joinMulticast(conn, opts.Interface, listenAddrWithPort)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return conn, nil
}

func handleGetDiscover(mgmtServerAddress string) mux.Handler {
	return mux.HandlerFunc(func(w mux.ResponseWriter, r *mux.Message) {
		log.Infof("Got discovery request from %v. Responding with management server address '%v'.",
			w.Client().RemoteAddr(), mgmtServerAddress,
		)
		coap_utils.RespondWithJSON(w, DiscoveryResponse{mgmtServerAddress})
	})
}

func joinMulticast(conn *net.UDPConn, ifaceStr string, addrStr string) error {
	iface, err := gonet.InterfaceByName(ifaceStr)
	if err != nil {
		return errors.WithStack(err)
	}

	addr, err := gonet.ResolveUDPAddr("udp", addrStr)
	if err != nil {
		return errors.WithStack(err)
	}

	err = conn.JoinGroup(iface, addr)
	if err != nil {
		return errors.WithStack(err)
	}

	log.Printf("Listening on %v with address %v", iface.Name, addr)

	return nil
}
