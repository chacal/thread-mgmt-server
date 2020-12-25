package http

import (
	"github.com/chacal/thread-mgmt-server/pkg/device_gateway"
	"github.com/chacal/thread-mgmt-server/pkg/device_registry"
	"github.com/chacal/thread-mgmt-server/pkg/state_poller_service"
	"github.com/gin-gonic/gin"
)

type depHandlerFunc = func(reg *device_registry.Registry, gw device_gateway.DeviceGateway, sps state_poller_service.StatePollerService, ctx *gin.Context)

func handlerWithDeps(reg *device_registry.Registry, gw device_gateway.DeviceGateway, sps state_poller_service.StatePollerService,
	f depHandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		f(reg, gw, sps, ctx)
	}
}

func handlerWithReg(reg *device_registry.Registry, f func(reg *device_registry.Registry, ctx *gin.Context)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		f(reg, ctx)
	}
}
