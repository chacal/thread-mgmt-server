package main

import (
	"context"
	"github.com/chacal/thread-mgmt-server/pkg/device_gateway"
	"github.com/chacal/thread-mgmt-server/pkg/device_registry"
	"github.com/chacal/thread-mgmt-server/pkg/server"
	"github.com/chacal/thread-mgmt-server/pkg/state_poller_service"
	log "github.com/sirupsen/logrus"
)

type Options struct {
	CoapPort int    `short:"c" long:"coap_port" description:"CoAP port to listen" default:"5683" env:"COAP_PORT"`
	HttpPort int    `short:"p" long:"http_port" description:"HTTP port to listen" default:"8080" env:"HTTP_PORT"`
	DbFile   string `short:"f" long:"file" description:"Database file for device registry" default:"devices.db" env:"DB_FILE"`
}

func main() {
	opts := Options{}
	server.ParseOptions(&opts)
	logOptions(opts)

	log.Info(Splash)

	// Start device registry
	reg, err := device_registry.Open(opts.DbFile)
	if err != nil {
		log.Fatalf("Failed to open device registry from file '%v'. Error: %+v", opts.DbFile, err)
	}
	defer reg.Close()

	gw := device_gateway.Create()

	sps := state_poller_service.Create(reg)
	err = sps.Refresh()
	if err != nil {
		log.Fatalf("Failed create state poller service. Error: %+v", err)
	}

	serverExit := make(chan int, 2)

	// Start CoAP server
	go startCoapServer(opts, reg, serverExit)

	// Start HTTP server
	go startHttpServer(opts, reg, gw, sps, serverExit)

	// Wait for servers to exit
	<-serverExit
	log.Fatalf("%+v", err)
}

func startCoapServer(opts Options, reg *device_registry.Registry, serverExit chan int) {
	coapServer, err := NewCoapServer(opts.CoapPort, reg)
	if err != nil {
		log.Fatalf("failed to create CoAP server: %+v", err)
	}
	defer coapServer.Stop()

	err = coapServer.Serve()
	serverExit <- 1
}

func startHttpServer(opts Options, reg *device_registry.Registry, gw device_gateway.DeviceGateway,
	sps state_poller_service.StatePollerService, serverExit chan int) {
	httpServer, err := NewHttpServer(opts, reg, gw, sps)
	if err != nil {
		log.Fatalf("failed to create HTTP server: %+v", err)
	}
	defer httpServer.Shutdown(context.Background())

	err = httpServer.ListenAndServe()
	serverExit <- 1
}

func logOptions(opts Options) {
	format := "Using configuration:\n" +
		"--------------------\n" +
		"CoAP listen port:\t%v\n" +
		"HTTP listen port:\t%v\n" +
		"DB file:\t\t%v\n"
	log.Infof(format, opts.CoapPort, opts.HttpPort, opts.DbFile)
}
