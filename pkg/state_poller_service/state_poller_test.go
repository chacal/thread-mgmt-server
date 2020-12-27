package state_poller_service

import (
	"github.com/chacal/thread-mgmt-server/pkg/device_gateway"
	"github.com/chacal/thread-mgmt-server/pkg/device_registry"
	"github.com/chacal/thread-mgmt-server/pkg/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net"
	"testing"
	"time"
)

var testState = device_registry.State{
	Addresses:  []net.IP{ip},
	Vcc:        2980,
	Instance:   "A100",
	TxPower:    0,
	PollPeriod: 1000,
	Parent:     device_registry.ParentInfo{Rloc16: "0x0400", LinkQualityIn: 3, LinkQualityOut: 2, AvgRssi: -75, LatestRssi: -72},
}

func TestStatePoller_Start(t *testing.T) {
	reg, mockGw := create(t)

	poller := createPoller(reg, mockGw, duration(t, "200ms"))
	defer poller.Stop()
	dev, _ := reg.Create("12345")
	assert.Equal(t, (*device_registry.State)(nil), dev.State)

	mockGw.EXPECT().FetchState(gomock.Eq(ip)).Return(testState, nil)
	poller.Start()

	// Wait for immediate poll
	time.Sleep(100 * time.Millisecond)

	dev, _ = reg.Get("12345")
	assert.Equal(t, &testState, dev.State)

	testState2 := testState
	testState2.Vcc = 3000
	mockGw.EXPECT().FetchState(gomock.Eq(ip)).Return(testState2, nil)

	// Wait for the first timer poll
	time.Sleep(200 * time.Millisecond)

	dev, _ = reg.Get("12345")
	assert.Equal(t, &testState2, dev.State)
}

func TestStatePoller_Refresh(t *testing.T) {
	reg, mockGw := create(t)

	poller := createPoller(reg, mockGw, duration(t, "200ms"))
	defer poller.Stop()
	_, _ = reg.Create("12345")

	mockGw.EXPECT().FetchState(gomock.Eq(ip)).Return(testState, nil)
	poller.Start()

	// Wait for immediate poll
	time.Sleep(100 * time.Millisecond)

	mockGw.EXPECT().FetchState(gomock.Eq(ip2)).Return(testState, nil)
	poller.Refresh(1, ip2)

	// Wait for the next poll
	time.Sleep(1100 * time.Millisecond)
}

func TestStatePoller_Stop(t *testing.T) {
	reg, mockGw := create(t)

	poller := createPoller(reg, mockGw, duration(t, "200ms"))
	_, _ = reg.Create("12345")

	mockGw.EXPECT().FetchState(gomock.Eq(ip)).Return(testState, nil)
	poller.Start()

	// Wait for immediate poll
	time.Sleep(100 * time.Millisecond)

	poller.Stop()

	// Wait for timer poll (shouldn't happen)
	time.Sleep(300 * time.Millisecond)
}

func create(t *testing.T) (*device_registry.Registry, *mocks.MockDeviceGateway) {
	reg := device_registry.CreateTestRegistry(t)
	mockCtrl := gomock.NewController(t)
	t.Cleanup(func() {
		mockCtrl.Finish()
	})
	mockGw := mocks.NewMockDeviceGateway(mockCtrl)
	return reg, mockGw
}

func createPoller(reg *device_registry.Registry, gw device_gateway.DeviceGateway, interval time.Duration) *statePoller {
	return &statePoller{"12345", interval, ip, nil,
		gw, reg, func() time.Duration { return 0 },
	}
}

func duration(t *testing.T, duration string) time.Duration {
	d, err := time.ParseDuration(duration)
	require.NoError(t, err)
	return d
}
