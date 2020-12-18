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
