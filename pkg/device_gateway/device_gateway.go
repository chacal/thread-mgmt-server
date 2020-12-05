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
	PushSettings(device device_registry.Device, destination net.IP) error
}

type deviceGateway struct{}

func Create() *deviceGateway {
	return &deviceGateway{}
}

func (r *deviceGateway) PushSettings(d device_registry.Device, destination net.IP) error {
	log.Debugf("Pushing settings %+v to %+v", d, destination)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	d.Addresses = []device_registry.DeviceAddress{} // Don't send IP addresses as device doesn't need them
	payload, err := json.Marshal(d)
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = coap_utils.PostJSON(ctx, destination.String()+":"+DEVICE_COAP_PORT, "api/settings", string(payload))
	return err
}
