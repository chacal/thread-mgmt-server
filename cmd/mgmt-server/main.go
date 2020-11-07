package main

import (
	"github.com/chacal/thread-mgmt-server/pkg/device_registry"
	"github.com/chacal/thread-mgmt-server/pkg/mgmt_routes"
	"github.com/chacal/thread-mgmt-server/pkg/server"
	"github.com/pkg/errors"
	"github.com/plgd-dev/go-coap/v2/mux"
	"github.com/plgd-dev/go-coap/v2/net"
	"github.com/plgd-dev/go-coap/v2/udp"
	log "github.com/sirupsen/logrus"
	"strconv"
)

type Options struct {
	Port   int    `short:"p" long:"port" description:"Port to listen" default:"5683" env:"PORT"`
	DbFile string `short:"f" long:"file" description:"Database file for device registry" default:"devices.db" env:"FILE"`
}

func main() {
	opts := Options{}
	server.ParseOptions(&opts)
	logOptions(opts)

	reg, err := device_registry.Open(opts.DbFile)
	if err != nil {
		log.Fatalf("Failed to open device registry from file '%v'. Error: %+v", opts.DbFile, err)
	}
	defer reg.Close()

	err = startCoapServer(opts, reg)
	log.Fatalf("%+v", err)
}

func startCoapServer(opts Options, reg *device_registry.Registry) error {
	conn, err := net.NewListenUDP("udp6", ":"+strconv.Itoa(opts.Port))
	if err != nil {
		return errors.WithStack(err)
	}
	defer conn.Close()

	router := mux.NewRouter()
	mgmt_routes.RegisterRoutes(router, reg)

	srv := udp.NewServer(udp.WithMux(router), udp.WithKeepAlive(nil))
	defer srv.Stop()
	return srv.Serve(conn)
}

func logOptions(opts Options) {
	format := "Using configuration:\n" +
		"--------------------\n" +
		"Listen port:\t%v\n"
	log.Infof(format, opts.Port)
}
