package state_poller_service

//go:generate mockgen -destination=../mocks/mock_state_poller_service.go -package=mocks github.com/chacal/thread-mgmt-server/pkg/state_poller_service StatePollerService

import (
	"github.com/chacal/thread-mgmt-server/pkg/device_registry"
	"net"
	"time"
)

type StatePollerService interface {
	Refresh() error
}

type statePollerService struct {
	reg           *device_registry.Registry
	pollers       map[string]StatePoller
	pollerCreator StatePollerCreator
}

func Create(reg *device_registry.Registry) *statePollerService {
	sp := statePollerService{reg, make(map[string]StatePoller), defaultStatePollerCreator}
	return &sp
}

func CreateWithPollerCreator(reg *device_registry.Registry, pollerCreator StatePollerCreator) *statePollerService {
	sp := statePollerService{reg, make(map[string]StatePoller), pollerCreator}
	return &sp
}

func (sp *statePollerService) Refresh() error {
	devices, err := sp.reg.GetDevices()
	if err != nil {
		return err
	}

	// Refresh existing devices
	for deviceId, device := range devices {
		poller, pollerExists := sp.pollers[deviceId]

		if device.Config.StatePollingEnabled && !pollerExists {
			sp.createPoller(deviceId, device.Config.StatePollingIntervalSec, device.Config.MainIp)
		} else if device.Config.StatePollingEnabled && pollerExists {
			poller.Refresh(device.Config.StatePollingIntervalSec, device.Config.MainIp)
		} else if !device.Config.StatePollingEnabled && pollerExists {
			sp.removePoller(deviceId)
		}
	}

	// Stop pollers for deleted devices
	for deviceId, _ := range sp.pollers {
		_, deviceExists := devices[deviceId]
		if !deviceExists {
			sp.removePoller(deviceId)
		}
	}

	return nil
}

func (sp *statePollerService) createPoller(deviceId string, pollingIntervalSec int, ip net.IP) {
	duration := time.Duration(pollingIntervalSec) * time.Second
	poller := sp.pollerCreator(sp.reg, deviceId, duration, ip)
	sp.pollers[deviceId] = poller
	poller.Start()
}

func (sp *statePollerService) removePoller(deviceId string) {
	poller := sp.pollers[deviceId]
	poller.Stop()
	delete(sp.pollers, deviceId)
}
