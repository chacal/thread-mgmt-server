package main

import (
	"github.com/chacal/thread-mgmt-server/pkg/device_registry"
	http_routes "github.com/chacal/thread-mgmt-server/pkg/mgmt_routes/http"
	T "github.com/chacal/thread-mgmt-server/pkg/test"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net"
	"testing"
)

var ip = net.ParseIP("ffff::1")

func TestV1GetDevices(t *testing.T) {
	router, reg := setup(t)

	T.AssertOKJson(t, `{}`, T.RecordGet(router, "/v1/devices"))

	err := reg.Update("12345", device_registry.Device{"D100", -4, 5000, nil})
	require.NoError(t, err)

	T.AssertOKJson(t, `{"12345": {"instance": "D100", "txPower": -4, "pollPeriod": 5000}}`, T.RecordGet(router, "/v1/devices"))

	err = reg.Update("ABCDE", device_registry.Device{"D100", -4, 5000, []net.IP{ip}})
	require.NoError(t, err)

	T.AssertOKJson(
		t,
		`{"12345": {"instance": "D100", "txPower": -4, "pollPeriod": 5000}, "ABCDE": {"instance": "D100", "txPower": -4, "pollPeriod": 5000, "addresses": ["ffff::1"]}}`,
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
			`{"id": "ABCDE", "instance": "D100", "txPower": -4, "pollPeriod": 5000, "addresses": ["ffff::1"]}`,
			device_registry.Device{"D100", -4, 5000, []net.IP{ip}},
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

func setup(t *testing.T) (*gin.Engine, *device_registry.Registry) {
	reg := device_registry.CreateTestRegistry(t)

	router := gin.Default()
	http_routes.RegisterRoutes(router, reg)

	return router, reg
}
