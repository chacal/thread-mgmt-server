package http

import (
	"github.com/chacal/thread-mgmt-server/pkg/device_registry"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func RegisterRoutes(router *gin.Engine, reg *device_registry.Registry) error {
	router.Use(errorHandlingMiddleware)
	router.Use(cors.Default())
	router.GET("/v1/devices", handlerWithReg(reg, getV1Devices))
	router.POST("/v1/devices/:device_id", handlerWithReg(reg, postV1Devices))
	router.DELETE("/v1/devices/:device_id", handlerWithReg(reg, deleteV1Devices))
	return serveStaticFromDir(router, "dist")
}

type Id struct {
	Id string `uri:"device_id" binding:"required"`
}

func getV1Devices(reg *device_registry.Registry, ctx *gin.Context) {
	devices, err := reg.GetAll()
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.IndentedJSON(http.StatusOK, devices)
}

func postV1Devices(reg *device_registry.Registry, ctx *gin.Context) {
	var id Id
	if err := ctx.ShouldBindUri(&id); err != nil {
		ctx.Error(errors.WithStack(err))
		return
	}

	var dev device_registry.Device
	if err := ctx.ShouldBindJSON(&dev); err != nil {
		ctx.Error(errors.WithStack(err))
		return
	}

	err := reg.Update(id.Id, dev)
	if err != nil {
		ctx.Error(errors.WithStack(err))
		return
	}

	ctx.Status(http.StatusOK)
}

func deleteV1Devices(reg *device_registry.Registry, ctx *gin.Context) {
	var id Id
	if err := ctx.ShouldBindUri(&id); err != nil {
		ctx.Error(errors.WithStack(err))
		return
	}

	err := reg.Delete(id.Id)
	if err != nil {
		ctx.Error(errors.WithStack(err))
		return
	}

	ctx.Status(http.StatusOK)
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

func serveStaticFromDir(router *gin.Engine, dir string) error {
	files, err := getFilenamesInDir(dir)
	if err != nil {
		return err
	}

	for _, f := range files {
		router.StaticFile(strings.TrimPrefix(f, dir), f)
	}

	// Manually add index.html
	router.StaticFile("/", dir + "/index.html")

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
