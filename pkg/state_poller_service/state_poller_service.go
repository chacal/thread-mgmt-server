package state_poller_service

//go:generate mockgen -destination=../mocks/mock_state_poller_service.go -package=mocks github.com/chacal/thread-mgmt-server/pkg/state_poller_service StatePollerService

import (
	"github.com/chacal/thread-mgmt-server/pkg/device_registry"
	"github.com/chacal/thread-mgmt-server/pkg/mqtt"
	log "github.com/sirupsen/logrus"
	"net"
	"time"
)

type pollResult struct {
	deviceId string
	state    device_registry.State
}

type StatePollerService interface {
	Start() error
	Stop()
	Refresh() error
}

type statePollerService struct {
	reg           *device_registry.Registry
	mqttSender    mqtt.MqttSender
	pollers       map[string]StatePoller
	pollerCreator StatePollerCreator
	pollResults   chan pollResult
	done          chan bool
}

func Create(reg *device_registry.Registry, mqttSender mqtt.MqttSender) *statePollerService {
	return CreateWithPollerCreator(reg, mqttSender, defaultStatePollerCreator)
}

func CreateWithPollerCreator(reg *device_registry.Registry, mqttSender mqtt.MqttSender, pollerCreator StatePollerCreator) *statePollerService {
	sp := statePollerService{
		reg:           reg,
		mqttSender:    mqttSender,
		pollers:       make(map[string]StatePoller),
		pollerCreator: pollerCreator,
		pollResults:   make(chan pollResult),
		done:          make(chan bool),
	}
	return &sp
}

func (sp *statePollerService) Start() error {
	go sp.handlePollResults()
	return sp.Refresh()
}

func (sp *statePollerService) Stop() {
	sp.done <- true
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

func (sp *statePollerService) handlePollResults() {
	for {
		select {
		case s := <-sp.pollResults:
			err := sp.reg.UpdateState(s.deviceId, s.state)
			if err != nil {
				log.Errorf("failed to update state, deviceId: %v, error: %v", s.deviceId, err)
			}
			sp.mqttSender.PublishState(s.state)
		case _ = <-sp.done:
			log.Infof("Ending poll result handling")
			return
		}
	}
}

func (sp *statePollerService) createPoller(deviceId string, pollingIntervalSec int, ip net.IP) {
	duration := time.Duration(pollingIntervalSec) * time.Second
	poller := sp.pollerCreator(sp.pollResults, deviceId, duration, ip)
	sp.pollers[deviceId] = poller
	poller.Start()
}

func (sp *statePollerService) removePoller(deviceId string) {
	poller := sp.pollers[deviceId]
	poller.Stop()
	delete(sp.pollers, deviceId)
}
