package state_poller_service

//go:generate mockgen -destination=../mocks/mock_state_poller.go -package=mocks github.com/chacal/thread-mgmt-server/pkg/state_poller_service StatePoller

import (
	"fmt"
	"github.com/chacal/thread-mgmt-server/pkg/device_gateway"
	"github.com/chacal/thread-mgmt-server/pkg/device_registry"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net"
	"time"
)

type StatePoller interface {
	Start()
	Refresh(pollingIntervalSec int, ip net.IP) error
	Stop()
}

type StatePollerCreator func(reg *device_registry.Registry, deviceId string, pollingInterval time.Duration, ip net.IP) StatePoller

type statePoller struct {
	deviceId             string
	statePollingInterval time.Duration
	ip                   net.IP
	ticker               *time.Ticker
	done                 chan bool
	gw                   device_gateway.DeviceGateway
	reg                  *device_registry.Registry
}

func defaultStatePollerCreator(reg *device_registry.Registry, deviceId string, pollingInterval time.Duration, ip net.IP) StatePoller {
	return &statePoller{deviceId, pollingInterval, ip,
		nil, make(chan bool), device_gateway.Create(), reg,
	}
}

func (sp *statePoller) Start() {
	log.Infof("Starting poller for device %v with interval %v", sp.deviceId, sp.statePollingInterval)
	sp.ticker = time.NewTicker(sp.statePollingInterval)
	go func() {
		sp.pollDeviceOnce()
		sp.pollWithTicker()
	}()
}

func (sp *statePoller) Refresh(pollingIntervalSec int, ip net.IP) error {
	duration, err := time.ParseDuration(fmt.Sprintf("%vs", pollingIntervalSec))
	if err != nil {
		return errors.WithStack(err)
	}
	if sp.statePollingInterval != duration || !sp.ip.Equal(ip) {
		log.Infof("Refreshing poller, interval: %v ip: %v", duration, ip)
	}
	return nil
}

func (sp *statePoller) Stop() {
	log.Infof("Stopping poller for device %v", sp.deviceId)
	sp.ticker.Stop()
	sp.done <- true
}

func (sp *statePoller) pollWithTicker() {
	for {
		select {
		case <-sp.done:
			log.Infof("Poller for %v done.", sp.deviceId)
			return
		case <-sp.ticker.C:
			sp.pollDeviceOnce()
		}
	}
}

func (sp *statePoller) pollDeviceOnce() {
	log.Infof("Polling device %v", sp.deviceId)
	state, err := sp.gw.FetchState(sp.ip)
	if err != nil {
		log.Errorf("failed to fetch state, deviceId: %v, ip: %v, error: %v", sp.deviceId, sp.ip, err)
	}

	log.Infof("Updating state for device %v: %v", sp.deviceId, state)
	err = sp.reg.UpdateState(sp.deviceId, state)
	if err != nil {
		log.Errorf("failed to update state, deviceId: %v, error: %v", sp.deviceId, err)
	}
}
