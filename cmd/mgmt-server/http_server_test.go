package main

import (
	"github.com/chacal/thread-mgmt-server/pkg/device_registry"
	http_routes "github.com/chacal/thread-mgmt-server/pkg/mgmt_routes/http"
	T "github.com/chacal/thread-mgmt-server/pkg/test"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestV1GetDevices(t *testing.T) {
	router, reg := setup(t)

	T.AssertOKJson(t, `{}`, T.RecordGet(router, "/v1/devices"))

	err := reg.Update("12345", device_registry.Device{"D100", 5000})
	require.NoError(t, err)

	T.AssertOKJson(t, `{"12345": {"name": "D100", "pollTime": 5000}}`, T.RecordGet(router, "/v1/devices"))
}

func TestV1PostDevice(t *testing.T) {
	router, reg := setup(t)

	T.AssertOK(t, T.RecordPost(router, "/v1/devices/12345", `{"id": "12345", "name": "D100", "pollTime": 5000}`))

	dev, err := reg.GetOrCreate("12345")
	require.NoError(t, err)
	assert.Equal(t, device_registry.Device{"D100", 5000}, dev)
}

func TestV1DeleteDevice(t *testing.T) {
	router, reg := setup(t)

	err := reg.Update("12345", device_registry.Device{"D100", 5000})
	require.NoError(t, err)

	T.AssertOK(t, T.RecordDelete(router, "/v1/devices/12345"))

	dev, err := reg.GetOrCreate("12345")
	require.NoError(t, err)
	assert.Equal(t, device_registry.Device{}, dev)  // Should return empty device
}

func setup(t *testing.T) (*gin.Engine, *device_registry.Registry) {
	reg := device_registry.CreateTestRegistry(t)

	router := gin.Default()
	http_routes.RegisterRoutes(router, reg)

	return router, reg
}
