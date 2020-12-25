package state_poller_service

//go:generate mockgen -destination=../mocks/mock_state_poller.go -package=mocks github.com/chacal/thread-mgmt-server/pkg/state_poller_service StatePoller

import (
	"fmt"
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

type StatePollerCreator func(deviceId string, pollingInterval time.Duration, ip net.IP) StatePoller

type statePoller struct {
	deviceId             string
	statePollingInterval time.Duration
	ip                   net.IP
}

func (sp *statePoller) Start() {
	log.Infof("Starting poller for device %v", sp.deviceId)
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
}
