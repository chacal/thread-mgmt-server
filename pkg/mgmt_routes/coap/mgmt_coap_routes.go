package coap

import (
	"encoding/json"
	"github.com/chacal/thread-mgmt-server/pkg/coap_utils"
	"github.com/chacal/thread-mgmt-server/pkg/device_registry"
	"github.com/pkg/errors"
	"github.com/plgd-dev/go-coap/v2/mux"
	"io/ioutil"
	"net"
)

func RegisterRoutes(router *mux.Router, reg *device_registry.Registry) {
	router.Use(coap_utils.LoggingMiddleware)
	router.Handle("v1/devices/", handlerWithReg(reg, getV1Device))
	router.Handle("v1/ip6/", handlerWithReg(reg, postV1IP6))
}

func getV1Device(reg *device_registry.Registry, w mux.ResponseWriter, r *mux.Message) {
	deviceId, err := coap_utils.GetLastPathPart(r)
	if err != nil {
		coap_utils.RespondWithInternalServerError(w, errors.WithStack(err))
	}

	dev, err := reg.GetOrCreate(deviceId)
	dev.Addresses = []device_registry.DeviceAddress{} // Don't send IP addresses as device doesn't need them
	if err != nil {
		coap_utils.RespondWithInternalServerError(w, errors.WithStack(err))
	}

	coap_utils.RespondWithJSON(w, dev)
}

func postV1IP6(reg *device_registry.Registry, w mux.ResponseWriter, r *mux.Message) {
	deviceId, err := coap_utils.GetLastPathPart(r)
	if err != nil {
		coap_utils.RespondWithInternalServerError(w, errors.WithStack(err))
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		coap_utils.RespondWithInternalServerError(w, errors.WithStack(err))
	}

	var addresses []net.IP
	err = json.Unmarshal(body, &addresses)
	if err != nil {
		coap_utils.RespondWithBadRequest(w, errors.WithStack(err))
	}

	err = reg.UpdateAddresses(deviceId, addresses)
	if err != nil {
		coap_utils.RespondWithInternalServerError(w, errors.WithStack(err))
	}

	coap_utils.RespondWithEmpty(w)
}

func handlerWithReg(reg *device_registry.Registry, f func(reg *device_registry.Registry, w mux.ResponseWriter, r *mux.Message)) mux.Handler {
	return mux.HandlerFunc(func(w mux.ResponseWriter, r *mux.Message) {
		f(reg, w, r)
	})
}
