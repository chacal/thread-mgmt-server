package coap

import (
	"encoding/json"
	"github.com/chacal/thread-mgmt-server/pkg/coap_utils"
	"github.com/chacal/thread-mgmt-server/pkg/device_registry"
	"github.com/pkg/errors"
	"github.com/plgd-dev/go-coap/v2/mux"
	"io/ioutil"
)

func RegisterRoutes(router *mux.Router, reg *device_registry.Registry) {
	router.Use(coap_utils.LoggingMiddleware)
	router.Handle("v1/defaults/", handlerWithReg(reg, getV1Defaults))
	router.Handle("v1/state/", handlerWithReg(reg, postV1State))
}

func getV1Defaults(reg *device_registry.Registry, w mux.ResponseWriter, r *mux.Message) {
	deviceId, err := coap_utils.GetLastPathPart(r)
	if err != nil {
		coap_utils.RespondWithInternalServerError(w, errors.WithStack(err))
	}

	dev, err := reg.GetOrCreate(deviceId)
	if err != nil {
		coap_utils.RespondWithInternalServerError(w, errors.WithStack(err))
	}

	coap_utils.RespondWithJSON(w, dev.Defaults)
}

func postV1State(reg *device_registry.Registry, w mux.ResponseWriter, r *mux.Message) {
	deviceId, err := coap_utils.GetLastPathPart(r)
	if err != nil {
		coap_utils.RespondWithInternalServerError(w, errors.WithStack(err))
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		coap_utils.RespondWithInternalServerError(w, errors.WithStack(err))
	}

	var state device_registry.State
	err = json.Unmarshal(body, &state)
	if err != nil {
		coap_utils.RespondWithBadRequest(w, errors.WithStack(err))
	}

	err = reg.UpdateState(deviceId, state)
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
