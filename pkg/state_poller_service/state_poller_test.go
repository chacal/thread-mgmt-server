package state_poller_service

import (
	"github.com/chacal/thread-mgmt-server/pkg/device_gateway"
	"github.com/chacal/thread-mgmt-server/pkg/device_registry"
	"github.com/chacal/thread-mgmt-server/pkg/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
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
	pollResults, mockGw := create(t)

	poller := createPoller(pollResults, mockGw, 200*time.Millisecond)
	defer poller.Stop()

	mockGw.EXPECT().FetchState(gomock.Eq(ip)).Return(testState, nil)
	poller.Start()

	// Wait for immediate poll
	result := <-pollResults
	assert.Equal(t, pollResult{"12345", testState}, result)

	testState2 := testState
	testState2.Vcc = 3000
	mockGw.EXPECT().FetchState(gomock.Eq(ip)).Return(testState2, nil)

	// Wait for the first timer poll
	result = <-pollResults
	assert.Equal(t, pollResult{"12345", testState2}, result)
}

func TestStatePoller_Refresh(t *testing.T) {
	pollResults, mockGw := create(t)

	poller := createPoller(pollResults, mockGw, 200*time.Millisecond)
	defer poller.Stop()

	mockGw.EXPECT().FetchState(gomock.Eq(ip)).Return(testState, nil)
	poller.Start()

	// Wait for immediate poll
	<-pollResults

	mockGw.EXPECT().FetchState(gomock.Eq(ip2)).Return(testState, nil)
	poller.Refresh(1, ip2)

	// Wait for the next poll
	<-pollResults
}

func TestStatePoller_Stop(t *testing.T) {
	pollResults, mockGw := create(t)

	poller := createPoller(pollResults, mockGw, 200*time.Millisecond)

	mockGw.EXPECT().FetchState(gomock.Eq(ip)).Return(testState, nil)
	poller.Start()

	// Wait for immediate poll
	<-pollResults

	poller.Stop()

	// Wait for timer poll (shouldn't happen)
	time.Sleep(300 * time.Millisecond)
}

func create(t *testing.T) (chan pollResult, *mocks.MockDeviceGateway) {
	pollResults := make(chan pollResult)
	mockCtrl := gomock.NewController(t)
	t.Cleanup(func() {
		mockCtrl.Finish()
	})
	mockGw := mocks.NewMockDeviceGateway(mockCtrl)
	return pollResults, mockGw
}

func createPoller(pollResults chan pollResult, gw device_gateway.DeviceGateway, interval time.Duration) *statePoller {
	return &statePoller{"12345", interval, ip, nil,
		gw, pollResults, func() time.Duration { return 0 },
	}
}
