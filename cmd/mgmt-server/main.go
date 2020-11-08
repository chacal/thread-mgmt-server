package main

import (
	"github.com/chacal/thread-mgmt-server/pkg/device_registry"
	"github.com/chacal/thread-mgmt-server/pkg/server"
	log "github.com/sirupsen/logrus"
)

type Options struct {
	CoapPort int    `short:"c" long:"coap_port" description:"CoAP port to listen" default:"5683" env:"COAP_PORT"`
	DbFile   string `short:"f" long:"file" description:"Database file for device registry" default:"devices.db" env:"DB_FILE"`
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

	coapServer, err := NewCoapServer(opts, reg)
	if err != nil {
		log.Fatalf("failed to create CoAP server: %+v", err)
	}
	defer coapServer.Stop()
	log.Fatalf("%+v", coapServer.Serve())
}

func logOptions(opts Options) {
	format := "Using configuration:\n" +
		"--------------------\n" +
		"CoAP listen port:\t%v\n"
	log.Infof(format, opts.CoapPort)
}
