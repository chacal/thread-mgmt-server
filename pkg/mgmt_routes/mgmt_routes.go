package mgmt_routes

import (
	"github.com/chacal/thread-mgmt-server/pkg/coap_utils"
	"github.com/chacal/thread-mgmt-server/pkg/device_registry"
	"github.com/pkg/errors"
	"github.com/plgd-dev/go-coap/v2/mux"
)

func RegisterRoutes(router *mux.Router, reg *device_registry.Registry) {
	router.Use(coap_utils.LoggingMiddleware)
	router.Handle("/v1/config/", handlerWithReg(reg, getV1Config))
}

func getV1Config(reg *device_registry.Registry, w mux.ResponseWriter, r *mux.Message) {
	deviceId, err := coap_utils.GetLastPathPart(r)
	if err != nil {
		coap_utils.RespondWithInternalServerError(w, errors.WithStack(err))
	}

	dev, err := reg.GetOrCreate(deviceId)
	if err != nil {
		coap_utils.RespondWithInternalServerError(w, errors.WithStack(err))
	}

	coap_utils.RespondWithJSON(w, dev)
}

func handlerWithReg(reg *device_registry.Registry, f func(reg *device_registry.Registry, w mux.ResponseWriter, r *mux.Message)) mux.Handler {
	return mux.HandlerFunc(func(w mux.ResponseWriter, r *mux.Message) {
		f(reg, w, r)
	})
}
