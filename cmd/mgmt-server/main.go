package main

import (
	"github.com/chacal/thread-mgmt-server/pkg/device_registry"
	"github.com/chacal/thread-mgmt-server/pkg/server"
	log "github.com/sirupsen/logrus"
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

	srv, err := NewServer(opts, reg)
	if err != nil {
		log.Fatalf("failed to open create server: %+v", err)
	}
	defer srv.Stop()
	log.Fatalf("%+v", srv.Serve())
}

func logOptions(opts Options) {
	format := "Using configuration:\n" +
		"--------------------\n" +
		"Listen port:\t%v\n"
	log.Infof(format, opts.Port)
}
