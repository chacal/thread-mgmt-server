package state_poller_service

//go:generate mockgen -destination=../mocks/mock_state_poller_service.go -package=mocks github.com/chacal/thread-mgmt-server/pkg/state_poller_service StatePollerService

import (
	"fmt"
	"github.com/chacal/thread-mgmt-server/pkg/device_registry"
	"github.com/pkg/errors"
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
			err = sp.createPoller(deviceId, device.Config.StatePollingIntervalSec, device.Config.MainIp)
			if err != nil {
				return nil
			}
		} else if device.Config.StatePollingEnabled && pollerExists {
			err = poller.Refresh(device.Config.StatePollingIntervalSec, device.Config.MainIp)
			if err != nil {
				return nil
			}
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

func (sp *statePollerService) createPoller(deviceId string, pollingIntervalSec int, ip net.IP) error {
	duration, err := time.ParseDuration(fmt.Sprintf("%vs", pollingIntervalSec))
	if err != nil {
		return errors.WithStack(err)
	}
	poller := sp.pollerCreator(deviceId, duration, ip)
	sp.pollers[deviceId] = poller
	poller.Start()
	return nil
}

func (sp *statePollerService) removePoller(deviceId string) {
	poller := sp.pollers[deviceId]
	poller.Stop()
	delete(sp.pollers, deviceId)
}

func defaultStatePollerCreator(deviceId string, pollingInterval time.Duration, ip net.IP) StatePoller {
	return &statePoller{deviceId, pollingInterval, ip}
}
