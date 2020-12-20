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
var ip2 = net.ParseIP("ffff::2")
var addr = []net.IP{ip}

func TestV1GetDevices(t *testing.T) {
	router, reg := setup(t)

	T.AssertOKJson(t, `{}`, T.RecordGet(router, "/v1/devices"))

	err := reg.UpdateDefaults("12345", device_registry.Defaults{"D100", T.IntP(-4), 5000})
	require.NoError(t, err)

	T.AssertOKJson(t,
		`{
			"12345": {
				"defaults": {
					"instance": "D100",
					"txPower": -4,
					"pollPeriod": 5000
				},
				"config": {
				},
				"state": {
				}
			}
		}`,
		T.RecordGet(router, "/v1/devices"),
	)

	err = reg.UpdateDefaults("ABCDE", device_registry.Defaults{"D100", T.IntP(-4), 5000})
	require.NoError(t, err)
	err = reg.UpdateState("ABCDE", DefaultState)
	require.NoError(t, err)

	T.AssertOKJson(t,
		`{
			"12345": {
				"defaults": { "instance": "D100", "txPower": -4, "pollPeriod": 5000 },
				"config": {},
				"state": {}
			},
			"ABCDE": {
				"defaults": { "instance": "D100", "txPower": -4, "pollPeriod": 5000 },
				"config": {},
				"state": {
					"vcc": 2970,
					"instance": "A100",
					"addresses": [
						"ffff::1"
					],
					"parent": {
						"rloc16": "0x4400",
						"linkQualityIn": 3,
						"linkQualityOut": 0,
						"avgRssi": -65,
						"latestRssi": -63
					}
				}
			}
		}`,
		T.RecordGet(router, "/v1/devices"),
	)
}

func TestV1PostDefaults(t *testing.T) {
	router, reg := setup(t)
	tests := map[string]struct {
		id       string
		payload  string
		expected device_registry.Defaults
	}{
		"default": {
			"12345",
			`{"instance": "D100", "txPower": -4, "pollPeriod": 5000}`,
			device_registry.Defaults{"D100", T.IntP(-4), 5000},
		},
		"missing poll period": {
			"ABCDE",
			`{"instance": "D100", "txPower": -4}`,
			device_registry.Defaults{Instance: "D100", TxPower: T.IntP(-4)},
		},
		"replaces previous defaults": {
			"ABCDE",
			`{"pollPeriod": 5000}`,
			device_registry.Defaults{PollPeriod: 5000},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			T.AssertOK(t, T.RecordPost(router, "/v1/devices/"+tc.id+"/defaults", tc.payload))
			defaults, err := reg.GetDefaults(tc.id)
			require.NoError(t, err)
			assert.Equal(t, &tc.expected, defaults)
		})
	}
}

func TestV1PostConfig(t *testing.T) {
	router, reg := setup(t)
	tests := map[string]struct {
		id       string
		payload  string
		expected device_registry.Config
	}{
		"default": {
			"12345",
			`{"mainIp": "ffff::1"}`,
			device_registry.Config{ip, nil, 0},
		},
		"replaces previous value": {
			"12345",
			`{"mainIp": "ffff::2", "statePollingEnabled": true, "statePollingIntervalSec": 300}`,
			device_registry.Config{ip2, T.BoolP(true), 300},
		},
		"empty config": {
			"12345",
			`{}`,
			device_registry.Config{},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			T.AssertOK(t, T.RecordPost(router, "/v1/devices/"+tc.id+"/config", tc.payload))
			devices, err := reg.GetDevices()
			require.NoError(t, err)
			config := devices[tc.id].Config
			require.NoError(t, err)
			assert.Equal(t, tc.expected, config)
		})
	}
}

func TestV1DeleteDevice(t *testing.T) {
	router, reg := setup(t)

	err := reg.UpdateDefaults("12345", device_registry.Defaults{"D100", T.IntP(-4), 5000})
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

	dev := device_registry.Defaults{"D100", T.IntP(-4), 5000}
	err := reg.UpdateDefaults("12345", dev)
	require.NoError(t, err)

	mockGw.EXPECT().PushDefaults(gomock.Eq(dev), gomock.Eq(net.ParseIP("ffff::1")))
	T.AssertOK(t, T.RecordPost(router, "/v1/devices/12345/push", `{"address": "ffff::1"}`))
}

func TestV1PostRefreshState(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockGw := mocks.NewMockDeviceGateway(mockCtrl)

	router, reg := setupWithGw(t, mockGw)
	T.AssertNotFound(t, T.RecordPost(router, "/v1/devices/12345/refresh_state", `{"address": "ffff::1"}`))
	T.AssertBadRequest(t, T.RecordPost(router, "/v1/devices/12345/refresh_state", ""))

	_, err := reg.GetOrCreate("12345")
	require.NoError(t, err)

	state := DefaultState
	mockGw.EXPECT().FetchState(gomock.Eq(net.ParseIP("ffff::1"))).Return(state, nil)

	T.AssertOKJson(t,
		`{
				"vcc": 2970,
				"instance": "A100",
				"addresses": [
					"ffff::1"
				],
				"parent": {
					"rloc16": "0x4400",
					"linkQualityIn": 3,
					"linkQualityOut": 0,
					"avgRssi": -65,
					"latestRssi": -63
				}
			}`,
		T.RecordPost(router, "/v1/devices/12345/refresh_state", `{"address": "ffff::1"}`),
	)

	devices, err := reg.GetDevices()
	require.NoError(t, err)
	assert.Equal(t, devices["12345"].State, state)
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
