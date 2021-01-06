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
var ip2 = net.ParseIP("ffff::2")

func TestCreate(t *testing.T) {
	reg := device_registry.CreateTestRegistry(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockSender := mocks.NewMockMqttSender(mockCtrl)
	sp := Create(reg, mockSender)

	assert.Empty(t, sp.pollers)
}

func TestStatePollerService_Refresh_whenConfigChanges(t *testing.T) {
	reg := device_registry.CreateTestRegistry(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockPoller := mocks.NewMockStatePoller(mockCtrl)
	mockSender := mocks.NewMockMqttSender(mockCtrl)
	sp := CreateWithPollerCreator(reg, mockSender, mockDevicePollerCreator(mockPoller))
	_, _ = reg.Create("12345")

	// Refresh with disabled polling should not start pollers
	reg.UpdateConfig("12345", device_registry.Config{ip, false, 600})
	err := sp.Refresh()
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
	mockSender := mocks.NewMockMqttSender(mockCtrl)
	sp := CreateWithPollerCreator(reg, mockSender, mockDevicePollerCreator(mockPoller))
	_, _ = reg.Create("12345")

	// Refresh with enabled polling should start poller
	reg.UpdateConfig("12345", device_registry.Config{ip, true, 600})
	mockPoller.EXPECT().Start()
	err := sp.Refresh()
	require.NoError(t, err)

	// Refresh with deleted device should stop poller
	reg.DeleteDevice("12345")
	mockPoller.EXPECT().Stop()
	err = sp.Refresh()
	require.NoError(t, err)
}

func TestStatePollerService_PollResultHandling(t *testing.T) {
	reg := device_registry.CreateTestRegistry(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockSender := mocks.NewMockMqttSender(mockCtrl)

	sps := Create(reg, mockSender)
	_, _ = reg.Create("12345")

	pollResults := make(chan pollResult)
	sps.pollResults = pollResults

	err := sps.Start()
	require.NoError(t, err)
	defer sps.Stop()

	mockSender.EXPECT().PublishState(gomock.Eq(testState))
	pollResults <- pollResult{"12345", testState}

	assert.Eventually(t, func() bool {
		dev, _ := reg.Get("12345")
		return assert.ObjectsAreEqual(&testState, dev.State)
	}, time.Second, 10*time.Millisecond)
}

func mockDevicePollerCreator(mockPoller *mocks.MockStatePoller) StatePollerCreator {
	return func(pollResults chan pollResult, deviceId string, pollingInterval time.Duration, ip net.IP) StatePoller {
		return mockPoller
	}
}
