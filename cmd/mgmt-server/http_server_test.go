package main

import (
	"github.com/chacal/thread-mgmt-server/pkg/device_gateway"
	"github.com/chacal/thread-mgmt-server/pkg/device_registry"
	http_routes "github.com/chacal/thread-mgmt-server/pkg/mgmt_routes/http"
	"github.com/chacal/thread-mgmt-server/pkg/mocks"
	"github.com/chacal/thread-mgmt-server/pkg/state_poller_service"
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

	_, err := reg.Create("12345")
	require.NoError(t, err)

	T.AssertOKJson(t,
		`{
			"12345": {
				"defaults": {
					"instance": "0000",
					"txPower": 0,
					"pollPeriod": 1000,
                    "displayType": ""
				},
				"config": {
					"mainIp": "",
					"statePollingEnabled": false,
					"statePollingIntervalSec": 600
				}
			}
		}`,
		T.RecordGet(router, "/v1/devices"),
	)

	err = reg.UpdateDefaults("12345",
		device_registry.Defaults{"D100", -4, 5000, device_registry.GOOD_DISPLAY_2_9IN_4GRAY},
	)
	require.NoError(t, err)

	T.AssertOKJson(t,
		`{
			"12345": {
				"defaults": {
					"instance": "D100",
					"txPower": -4,
					"pollPeriod": 5000,
					"displayType": "GOOD_DISPLAY_2_9IN_4GRAY"
				},
				"config": {
					"mainIp": "",
					"statePollingEnabled": false,
					"statePollingIntervalSec": 600
				}
			}
		}`,
		T.RecordGet(router, "/v1/devices"),
	)

	_, err = reg.Create("ABCDE")
	require.NoError(t, err)

	err = reg.UpdateDefaults("ABCDE",
		device_registry.Defaults{"D101", 0, 3000, device_registry.GOOD_DISPLAY_2_9IN},
	)
	require.NoError(t, err)
	err = reg.UpdateState("ABCDE", testState)
	require.NoError(t, err)

	T.AssertOKJson(t,
		`{
			"12345": {
				"defaults": { "instance": "D100", "txPower": -4, "pollPeriod": 5000, "displayType": "GOOD_DISPLAY_2_9IN_4GRAY" },
				"config": { "mainIp": "", "statePollingEnabled": false, "statePollingIntervalSec": 600 }
			},
			"ABCDE": {
				"defaults": { "instance": "D101", "txPower": 0, "pollPeriod": 3000, "displayType": "GOOD_DISPLAY_2_9IN" },
				"config": { "mainIp": "", "statePollingEnabled": false, "statePollingIntervalSec": 600 },
				"state": {
					"vcc": 2970,
					"instance": "A100",
					"addresses": [
						"ffff::1"
					],
					"txPower": -4,
					"pollPeriod": 1000,
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
			"ABCDE",
			`{"instance": "D101", "txPower": 0, "pollPeriod": 1000, "displayType": ""}`,
			device_registry.Defaults{"D101", 0, 1000, ""},
		},
		"replaces previous defaults": {
			"ABCDE",
			`{"instance": "D102", "txPower": 4, "pollPeriod": 2000, "displayType": "GOOD_DISPLAY_1_54IN"}`,
			device_registry.Defaults{"D102", 4, 2000, device_registry.GOOD_DISPLAY_1_54IN},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			_, _ = reg.Create(tc.id)
			T.AssertOK(t, T.RecordPost(router, "/v1/devices/"+tc.id+"/defaults", tc.payload))
			device, err := reg.Get(tc.id)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, device.Defaults)
		})
	}
}

func TestV1PostConfig(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockSps := mocks.NewMockStatePollerService(mockCtrl)

	router, reg := setupWithSps(t, mockSps)
	tests := map[string]struct {
		id       string
		payload  string
		expected device_registry.Config
	}{
		"default": {
			"12345",
			`{"mainIp": "ffff::1", "statePollingEnabled": false, "statePollingIntervalSec": 100}`,
			device_registry.Config{ip, false, 100},
		},
		"replaces previous value": {
			"12345",
			`{"mainIp": "ffff::2", "statePollingEnabled": true, "statePollingIntervalSec": 300}`,
			device_registry.Config{ip2, true, 300},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			_, _ = reg.Create(tc.id)
			mockSps.EXPECT().Refresh()
			T.AssertOK(t, T.RecordPost(router, "/v1/devices/"+tc.id+"/config", tc.payload))
			device, err := reg.Get(tc.id)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, device.Config)
		})
	}
}

func TestV1DeleteDevice(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockSps := mocks.NewMockStatePollerService(mockCtrl)

	router, reg := setupWithSps(t, mockSps)

	_, err := reg.Create("12345")
	require.NoError(t, err)

	mockSps.EXPECT().Refresh()

	T.AssertOK(t, T.RecordDelete(router, "/v1/devices/12345"))

	contains, err := reg.Contains("12345")
	require.NoError(t, err)
	assert.Equal(t, false, contains)
}

func TestV1PostDevicePush(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockGw := mocks.NewMockDeviceGateway(mockCtrl)

	router, reg := setupWithGw(t, mockGw)
	T.AssertNotFound(t, T.RecordPost(router, "/v1/devices/12345/push", `{"address": "ffff::1"}`))
	T.AssertBadRequest(t, T.RecordPost(router, "/v1/devices/12345/push", ""))

	_, err := reg.Create("12345")
	require.NoError(t, err)

	mockGw.EXPECT().PushDefaults(gomock.Eq(device_registry.DefaultDevice.Defaults), gomock.Eq(net.ParseIP("ffff::1")))
	T.AssertOK(t, T.RecordPost(router, "/v1/devices/12345/push", `{"address": "ffff::1"}`))
}

func TestV1PostRefreshState(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockGw := mocks.NewMockDeviceGateway(mockCtrl)

	router, reg := setupWithGw(t, mockGw)
	T.AssertNotFound(t, T.RecordPost(router, "/v1/devices/12345/refresh_state", `{"address": "ffff::1"}`))
	T.AssertBadRequest(t, T.RecordPost(router, "/v1/devices/12345/refresh_state", ""))

	_, err := reg.Create("12345")
	require.NoError(t, err)

	state := testState
	mockGw.EXPECT().FetchState(gomock.Eq(net.ParseIP("ffff::1"))).Return(state, nil)

	T.AssertOKJson(t,
		`{
				"vcc": 2970,
				"instance": "A100",
				"addresses": [
					"ffff::1"
				],
				"txPower": -4,
				"pollPeriod": 1000,
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

	device, err := reg.Get("12345")
	require.NoError(t, err)
	assert.Equal(t, *device.State, state)
}

func setup(t *testing.T) (*gin.Engine, *device_registry.Registry) {
	gw := device_gateway.Create()
	return setupWithGw(t, gw)
}

func setupWithGw(t *testing.T, gw device_gateway.DeviceGateway) (*gin.Engine, *device_registry.Registry) {
	reg := device_registry.CreateTestRegistry(t)
	sps := state_poller_service.Create(reg)
	router := gin.Default()
	http_routes.RegisterRoutes(router, reg, gw, sps)

	return router, reg
}

func setupWithSps(t *testing.T, sps state_poller_service.StatePollerService) (*gin.Engine, *device_registry.Registry) {
	reg := device_registry.CreateTestRegistry(t)
	gw := device_gateway.Create()
	router := gin.Default()
	http_routes.RegisterRoutes(router, reg, gw, sps)

	return router, reg
}
