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
	dbFile := T.Tempfile()
	reg, err := device_registry.Open(dbFile)
	require.NoError(t, err)
	defer reg.Close()

	router := gin.Default()
	http_routes.RegisterRoutes(router, reg)

	T.AssertOKJson(t, `{}`, T.RecordGet(router, "/v1/devices"))

	err = reg.Update("12345", device_registry.Device{"D100", 5000})
	require.NoError(t, err)

	T.AssertOKJson(t, `{"12345": {"name": "D100", "pollTime": 5000}}`, T.RecordGet(router, "/v1/devices"))
}

func TestV1PostDevice(t *testing.T) {
	dbFile := T.Tempfile()
	reg, err := device_registry.Open(dbFile)
	require.NoError(t, err)
	defer reg.Close()

	router := gin.Default()
	http_routes.RegisterRoutes(router, reg)

	T.AssertOK(t, T.RecordPost(router, "/v1/devices/12345", `{"id": "12345", "name": "D100", "pollTime": 5000}`))

	dev, err := reg.GetOrCreate("12345")
	require.NoError(t, err)
	assert.Equal(t, device_registry.Device{"D100", 5000}, dev)
}
