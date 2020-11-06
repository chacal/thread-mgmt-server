package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	gonet "net"
	"os"
)

func ParseOptions() Options {
	opts := Options{}
	var parser = flags.NewParser(&opts, flags.Default)
	if _, err := parser.Parse(); err != nil {
		flagsErr, ok := err.(*flags.Error)
		if ok && flagsErr.Type != flags.ErrHelp {
			fmt.Println()
			parser.WriteHelp(os.Stdout)
		}
		os.Exit(1)
	}
	return opts
}

func LogOptions(opts Options) {
	format := "Using configuration:\n" +
		"--------------------\n" +
		"Interface:\t\t%v\n" +
		"Listen address:\t\t%v\n" +
		"Listen port:\t\t%v\n" +
		"Mgmt server address:\t%v\n"
	log.Infof(format, opts.Interface, opts.ListenAddr, opts.Port, opts.MgmtServerAddress)
}

func findLoopbackInterface() (*gonet.Interface, error) {
	ifaces, err := gonet.Interfaces()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	for _, i := range ifaces {
		if i.Flags&gonet.FlagLoopback > 0 {
			return &i, nil
		}
	}
	return nil, errors.Errorf("Couldn't find loopback interface! Interfaces: %v", ifaces)
}
