package main

import (
	"github.com/chacal/thread-mgmt-server/pkg/device_gateway"
	"github.com/chacal/thread-mgmt-server/pkg/device_registry"
	http_routes "github.com/chacal/thread-mgmt-server/pkg/mgmt_routes/http"
	"github.com/chacal/thread-mgmt-server/pkg/mocks"
	T "github.com/chacal/thread-mgmt-server/pkg/test"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net"
	"testing"
)

var ip = net.ParseIP("ffff::1")
var addr = []device_registry.DeviceAddress{{ip, false}}

func TestV1GetDevices(t *testing.T) {
	router, reg := setup(t)

	T.AssertOKJson(t, `{}`, T.RecordGet(router, "/v1/devices"))

	err := reg.Update("12345", device_registry.Device{"D100", -4, 5000, nil})
	require.NoError(t, err)

	T.AssertOKJson(t, `{"12345": {"instance": "D100", "txPower": -4, "pollPeriod": 5000}}`, T.RecordGet(router, "/v1/devices"))

	err = reg.Update("ABCDE", device_registry.Device{"D100", -4, 5000, addr})
	require.NoError(t, err)

	T.AssertOKJson(
		t,
		`{"12345": {"instance": "D100", "txPower": -4, "pollPeriod": 5000}, "ABCDE": {"instance": "D100", "txPower": -4, "pollPeriod": 5000, "addresses": [{"ip": "ffff::1", "main": false}]}}`,
		T.RecordGet(router, "/v1/devices"),
	)
}

func TestV1PostDevice(t *testing.T) {
	router, reg := setup(t)
	tests := map[string]struct {
		id       string
		payload  string
		expected device_registry.Device
	}{
		"no addresses": {
			"12345",
			`{"id": "12345", "instance": "D100", "txPower": -4, "pollPeriod": 5000}`,
			device_registry.Device{"D100", -4, 5000, nil},
		},
		"empty addresses": {
			"ABCDE",
			`{"id": "ABCDE", "instance": "D100", "txPower": -4, "pollPeriod": 5000, "addresses": []}`,
			device_registry.Device{"D100", -4, 5000, nil},
		},
		"with address": {
			"ABCDE",
			`{"id": "ABCDE", "instance": "D100", "txPower": -4, "pollPeriod": 5000, "addresses": [{"ip": "ffff::1", "main": false}]}`,
			device_registry.Device{"D100", -4, 5000, addr},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			T.AssertOK(t, T.RecordPost(router, "/v1/devices/"+tc.id, tc.payload))
			dev, err := reg.GetOrCreate(tc.id)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, dev)
		})
	}
}

func TestV1DeleteDevice(t *testing.T) {
	router, reg := setup(t)

	err := reg.Update("12345", device_registry.Device{"D100", -4, 5000, nil})
	require.NoError(t, err)

	T.AssertOK(t, T.RecordDelete(router, "/v1/devices/12345"))

	dev, err := reg.GetOrCreate("12345")
	require.NoError(t, err)
	assert.Equal(t, device_registry.Device{}, dev) // Should return empty device
}

func TestV1PostDevicePush(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockGw := mocks.NewMockDeviceGateway(mockCtrl)

	router, reg := setupWithGw(t, mockGw)
	T.AssertNotFound(t, T.RecordPost(router, "/v1/devices/12345/push", `{"address": "ffff::1"}`))
	T.AssertBadRequest(t, T.RecordPost(router, "/v1/devices/12345/push", ""))

	dev := device_registry.Device{"D100", -4, 5000, nil}
	err := reg.Update("12345", dev)
	require.NoError(t, err)

	mockGw.EXPECT().PushSettings(gomock.Eq(dev), gomock.Eq(net.ParseIP("ffff::1")))
	T.AssertOK(t, T.RecordPost(router, "/v1/devices/12345/push", `{"address": "ffff::1"}`))
}

func setup(t *testing.T) (*gin.Engine, *device_registry.Registry) {
	gw := device_gateway.Create()
	return setupWithGw(t, gw)
}

func setupWithGw(t *testing.T, gw device_gateway.DeviceGateway) (*gin.Engine, *device_registry.Registry) {
	reg := device_registry.CreateTestRegistry(t)
	router := gin.Default()
	http_routes.RegisterRoutes(router, reg, gw)

	return router, reg
}
