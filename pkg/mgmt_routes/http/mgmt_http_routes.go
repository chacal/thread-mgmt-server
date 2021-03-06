package http

import (
	"github.com/chacal/thread-mgmt-server/pkg/device_gateway"
	"github.com/chacal/thread-mgmt-server/pkg/device_registry"
	"github.com/chacal/thread-mgmt-server/pkg/state_poller_service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func RegisterRoutes(router *gin.Engine, reg *device_registry.Registry, gw device_gateway.DeviceGateway,
	sps state_poller_service.StatePollerService) error {
	router.Use(errorHandlingMiddleware)
	router.Use(cors.Default())
	router.GET("/v1/devices", handlerWithReg(reg, getV1Devices))
	router.POST("/v1/devices/:device_id/defaults", handlerWithReg(reg, postV1Defaults))
	router.POST("/v1/devices/:device_id/config", handlerWithDeps(reg, gw, sps, postV1Config))
	router.POST("/v1/devices/:device_id/push", handlerWithDeps(reg, gw, sps, postV1DevicesPushDefaults))
	router.POST("/v1/devices/:device_id/refresh_state", handlerWithDeps(reg, gw, sps, postV1DevicesRefreshState))
	router.DELETE("/v1/devices/:device_id", handlerWithDeps(reg, gw, sps, deleteV1Device))
	return serveStaticFromDir(router, "dist")
}

type Id struct {
	Id string `uri:"device_id" binding:"required"`
}

func getV1Devices(reg *device_registry.Registry, ctx *gin.Context) {
	devices, err := reg.GetDevices()
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.IndentedJSON(http.StatusOK, devices)
}

func postV1Defaults(reg *device_registry.Registry, ctx *gin.Context) {
	var id Id
	if err := ctx.ShouldBindUri(&id); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, errors.WithStack(err))
		return
	}

	var defaults device_registry.Defaults
	if err := ctx.ShouldBindJSON(&defaults); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, errors.WithStack(err))
		return
	}

	err := reg.UpdateDefaults(id.Id, defaults)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusOK)
}

func postV1Config(reg *device_registry.Registry, gw device_gateway.DeviceGateway, sps state_poller_service.StatePollerService, ctx *gin.Context) {
	var id Id
	if err := ctx.ShouldBindUri(&id); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, errors.WithStack(err))
		return
	}

	var config device_registry.Config
	if err := ctx.ShouldBindJSON(&config); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, errors.WithStack(err))
		return
	}

	err := reg.UpdateConfig(id.Id, config)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = sps.Refresh()
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusOK)
}

func deleteV1Device(reg *device_registry.Registry, gw device_gateway.DeviceGateway, sps state_poller_service.StatePollerService, ctx *gin.Context) {
	var id Id
	if err := ctx.ShouldBindUri(&id); err != nil {
		ctx.Error(errors.WithStack(err))
		return
	}

	err := reg.DeleteDevice(id.Id)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = sps.Refresh()
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusOK)
}

type DeviceDestination struct {
	Address net.IP `json:"address" binding:"required"`
}

func postV1DevicesPushDefaults(reg *device_registry.Registry, gw device_gateway.DeviceGateway, sps state_poller_service.StatePollerService, ctx *gin.Context) {
	id, dst, err := assertDeviceFromRequestExists(reg, ctx)
	if err != nil {
		return
	}

	device, err := reg.Get(id)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = gw.PushDefaults(device.Defaults, dst)
	if err != nil {
		ctx.Error(errors.WithStack(err))
		return
	}

	ctx.Status(http.StatusOK)
}

func postV1DevicesRefreshState(reg *device_registry.Registry, gw device_gateway.DeviceGateway, sps state_poller_service.StatePollerService, ctx *gin.Context) {
	id, dst, err := assertDeviceFromRequestExists(reg, ctx)
	if err != nil {
		return
	}

	state, err := gw.FetchState(dst)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = reg.UpdateState(id, state)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.IndentedJSON(http.StatusOK, state)
}

func assertDeviceFromRequestExists(reg *device_registry.Registry, ctx *gin.Context) (string, net.IP, error) {
	var id Id
	if err := ctx.ShouldBindUri(&id); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, errors.WithStack(err))
		return "", nil, err
	}

	var dst DeviceDestination
	if err := ctx.ShouldBindJSON(&dst); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, errors.WithStack(err))
		return "", nil, err
	}

	deviceExists, err := reg.Contains(id.Id)
	if err != nil {
		ctx.Error(err)
		return "", nil, err
	}

	if !deviceExists {
		ctx.AbortWithStatus(http.StatusNotFound)
		return "", nil, errors.Errorf("device with id %v not found", id.Id)
	}
	return id.Id, dst.Address, nil
}


func errorHandlingMiddleware(ctx *gin.Context) {
	ctx.Next()
	if len(ctx.Errors) > 0 {
		for _, e := range ctx.Errors {
			log.Errorf("%+v", e.Err)
		}
		if !ctx.IsAborted() {
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
	}
}

func serveStaticFromDir(router *gin.Engine, dir string) error {
	files, err := getFilenamesInDir(dir)
	if err != nil {
		return err
	}

	for _, f := range files {
		router.StaticFile(strings.TrimPrefix(f, dir), f)
	}

	// Manually add index.html
	router.StaticFile("/", dir+"/index.html")

	return nil
}

func getFilenamesInDir(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info != nil && !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, errors.WithStack(err)
}
