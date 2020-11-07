package main

import (
	"fmt"
	"github.com/chacal/thread-mgmt-server/pkg/coap-utils"
	"github.com/pkg/errors"
	"github.com/plgd-dev/go-coap/v2/message"
	"github.com/plgd-dev/go-coap/v2/message/codes"
	"github.com/plgd-dev/go-coap/v2/mux"
	"github.com/plgd-dev/go-coap/v2/net"
	"github.com/plgd-dev/go-coap/v2/udp"
	log "github.com/sirupsen/logrus"
	gonet "net"
	"strconv"
	"strings"
)

type Options struct {
	Interface         string `short:"i" long:"interface" description:"Interface to bind to" default:"wpan0" env:"BIND_INTERFACE"`
	ListenAddr        string `short:"l" long:"listen" description:"Multicast address to listen" default:"[ff03::1]" env:"LISTEN_ADDRESS"`
	Port              int    `short:"p" long:"port" description:"Port to listen" default:"5683" env:"PORT"`
	MgmtServerAddress string `short:"m" long:"mgmt-server" description:"Address of the returned management server" env:"MGMT_SERVER_ADDRESS" required:"yes"`
}

func main() {
	opts := Options{}
	server.ParseOptions(opts)
	LogOptions(opts)
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
		payload := fmt.Sprintf(`{"mgmtServer": "%v"}`, mgmtServerAddress)
		log.Infof("Got discovery request from %v. Responding with management server address '%v'.",
			w.Client().RemoteAddr(), mgmtServerAddress,
		)
		w.SetResponse(codes.Content, message.AppJSON, strings.NewReader(payload))
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
