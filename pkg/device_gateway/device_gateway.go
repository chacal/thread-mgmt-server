package device_gateway

//go:generate mockgen -destination=../mocks/mock_device_gateway.go -package=mocks github.com/chacal/thread-mgmt-server/pkg/device_gateway DeviceGateway

import (
	"context"
	"encoding/json"
	"github.com/chacal/thread-mgmt-server/pkg/coap_utils"
	"github.com/chacal/thread-mgmt-server/pkg/device_registry"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net"
	"time"
)

var DEVICE_COAP_PORT = "5683"

type DeviceGateway interface {
	PushDefaults(defaults device_registry.Defaults, destination net.IP) error
	FetchState(destination net.IP) (device_registry.State, error)
}

type deviceGateway struct{}

func Create() *deviceGateway {
	return &deviceGateway{}
}

func (r *deviceGateway) PushDefaults(defaults device_registry.Defaults, destination net.IP) error {
	log.Debugf("Pushing settings %+v to %+v", defaults, destination)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	payload, err := json.Marshal(defaults)
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = coap_utils.PostJSON(ctx, "["+destination.String()+"]:"+DEVICE_COAP_PORT, "api/settings", string(payload))
	return err
}

func (r *deviceGateway) FetchState(destination net.IP) (device_registry.State, error) {
	log.Debugf("Fetching state from %+v", destination)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	res, err := coap_utils.GetJSON(ctx, "["+destination.String()+"]:"+DEVICE_COAP_PORT, "api/state")
	if err != nil {
		return device_registry.State{}, err
	}
	return device_registry.StateFromJSON([]byte(res))
}
