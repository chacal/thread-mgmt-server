package state_poller_service

//go:generate mockgen -destination=../mocks/mock_state_poller.go -package=mocks github.com/chacal/thread-mgmt-server/pkg/state_poller_service StatePoller

import (
	"github.com/chacal/thread-mgmt-server/pkg/device_gateway"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"net"
	"time"
)

var maxSleepRandomnessSeconds = 60

type StatePoller interface {
	Start()
	Refresh(pollingIntervalSec int, ip net.IP)
	Stop()
}

type StatePollerCreator func(pollResults chan pollResult, deviceId string, pollingInterval time.Duration, ip net.IP) StatePoller

type statePoller struct {
	deviceId             string
	statePollingInterval time.Duration
	ip                   net.IP
	timer                *time.Timer
	gw                   device_gateway.DeviceGateway
	pollResults          chan pollResult
	sleepRandomizer      func() time.Duration
}

func defaultStatePollerCreator(pollResults chan pollResult, deviceId string, pollingInterval time.Duration, ip net.IP) StatePoller {
	return &statePoller{deviceId, pollingInterval, ip, nil, device_gateway.Create(),
		pollResults, nextSleepRandomDuration,
	}
}

func (sp *statePoller) Start() {
	initialSleep := sp.sleepRandomizer()
	log.Infof("Starting poller for device %v with interval %v and initial sleep %v",
		sp.deviceId, sp.statePollingInterval, initialSleep,
	)
	sp.timer = time.AfterFunc(initialSleep, sp.pollDeviceOnce)
}

func (sp *statePoller) Refresh(pollingIntervalSec int, ip net.IP) {
	duration := time.Duration(pollingIntervalSec) * time.Second
	if sp.statePollingInterval != duration || !sp.ip.Equal(ip) {
		log.Infof("Refreshing poller, interval: %v ip: %v", duration, ip)
		// Stop & drain timer
		if !sp.timer.Stop() {
			<-sp.timer.C
		}
		sp.statePollingInterval = duration
		sp.timer.Reset(sp.statePollingInterval)
		sp.ip = ip
	}
}

func (sp *statePoller) Stop() {
	log.Infof("Stopping poller for device %v", sp.deviceId)
	sp.timer.Stop()
}

func (sp *statePoller) pollDeviceOnce() {
	var nextSleep = sp.statePollingInterval + sp.sleepRandomizer()
	log.Debugf("Polling device %v, next sleep %v", sp.deviceId, nextSleep)
	defer sp.timer.Reset(nextSleep)

	state, err := sp.gw.FetchState(sp.ip)
	if err != nil {
		log.Errorf("failed to fetch state, deviceId: %v, ip: %v, error: %v", sp.deviceId, sp.ip, err)
		return
	}

	sp.pollResults <- pollResult{sp.deviceId, state}
}

func nextSleepRandomDuration() time.Duration {
	return time.Duration(rand.Intn(maxSleepRandomnessSeconds)) * time.Second
}
