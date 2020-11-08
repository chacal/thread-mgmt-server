package http

import (
	"github.com/chacal/thread-mgmt-server/pkg/device_registry"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func RegisterRoutes(router *gin.Engine, reg *device_registry.Registry) {
	router.Use(errorHandlingMiddleware)
	router.GET("/v1/devices", handlerWithReg(reg, getV1Devices))
}

func getV1Devices(reg *device_registry.Registry, ctx *gin.Context) {
	devices, err := reg.GetAll()
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.IndentedJSON(http.StatusOK, devices)
}

func handlerWithReg(reg *device_registry.Registry, f func(reg *device_registry.Registry, ctx *gin.Context)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		f(reg, ctx)
	}
}

func errorHandlingMiddleware(ctx *gin.Context) {
	ctx.Next()
	if len(ctx.Errors) > 0 {
		for _, e := range ctx.Errors {
			log.Errorf("%+v", e.Err)
		}
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}
}
