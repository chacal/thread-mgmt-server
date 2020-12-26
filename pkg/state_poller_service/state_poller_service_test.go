package state_poller_service

import (
	"github.com/chacal/thread-mgmt-server/pkg/device_registry"
	"github.com/chacal/thread-mgmt-server/pkg/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net"
	"testing"
	"time"
)

var ip = net.ParseIP("ffff::1")
var ip2 = net.ParseIP("ffff::1")

func TestCreate(t *testing.T) {
	reg := device_registry.CreateTestRegistry(t)
	sp := Create(reg)

	assert.Empty(t, sp.pollers)
}

func TestStatePollerService_Refresh_whenConfigChanges(t *testing.T) {
	reg := device_registry.CreateTestRegistry(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockPoller := mocks.NewMockStatePoller(mockCtrl)
	sp := CreateWithPollerCreator(reg, mockDevicePollerCreator(mockPoller))
	_, err := reg.Create("12345")

	// Refresh with disabled polling should not start pollers
	reg.UpdateConfig("12345", device_registry.Config{ip, false, 600})
	err = sp.Refresh()
	require.NoError(t, err)

	// Refresh with enabled polling should start poller
	reg.UpdateConfig("12345", device_registry.Config{ip, true, 600})
	mockPoller.EXPECT().Start()
	err = sp.Refresh()
	require.NoError(t, err)

	// Refresh with changed config should refresh poller with the new config
	reg.UpdateConfig("12345", device_registry.Config{ip2, true, 500})
	mockPoller.EXPECT().Refresh(gomock.Eq(500), gomock.Eq(ip2))
	err = sp.Refresh()
	require.NoError(t, err)

	// Refresh with the same config should refresh poller with the same config
	mockPoller.EXPECT().Refresh(gomock.Eq(500), gomock.Eq(ip2))
	err = sp.Refresh()
	require.NoError(t, err)

	// Refresh with disabled polling should stop poller
	reg.UpdateConfig("12345", device_registry.Config{ip2, false, 500})
	mockPoller.EXPECT().Stop()
	err = sp.Refresh()
	require.NoError(t, err)

	// Refresh with no pollers should not do anything
	err = sp.Refresh()
	require.NoError(t, err)
}

func TestStatePollerService_Refresh_whenDeviceDeleted(t *testing.T) {
	reg := device_registry.CreateTestRegistry(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockPoller := mocks.NewMockStatePoller(mockCtrl)
	sp := CreateWithPollerCreator(reg, mockDevicePollerCreator(mockPoller))
	_, err := reg.Create("12345")

	// Refresh with enabled polling should start poller
	reg.UpdateConfig("12345", device_registry.Config{ip, true, 600})
	mockPoller.EXPECT().Start()
	err = sp.Refresh()
	require.NoError(t, err)

	// Refresh with deleted device should stop poller
	reg.DeleteDevice("12345")
	mockPoller.EXPECT().Stop()
	err = sp.Refresh()
	require.NoError(t, err)
}

func mockDevicePollerCreator(mockPoller *mocks.MockStatePoller) StatePollerCreator {
	return func(reg *device_registry.Registry, deviceId string, pollingInterval time.Duration, ip net.IP) StatePoller {
		return mockPoller
	}
}
