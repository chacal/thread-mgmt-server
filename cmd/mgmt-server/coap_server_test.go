package main

import (
	"context"
	"github.com/chacal/thread-mgmt-server/pkg/coap_utils"
	"github.com/chacal/thread-mgmt-server/pkg/device_registry"
	"github.com/plgd-dev/go-coap/v2/mux"
	"github.com/plgd-dev/go-coap/v2/udp/message/pool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net"
	"strconv"
	"testing"
	"time"
)

const TEST_COAP_PORT = 55683

var TEST_COAP_URL = "localhost:" + strconv.Itoa(TEST_COAP_PORT)

func TestGetV1Device(t *testing.T) {
	coapServerTest(t, func(t *testing.T, reg *device_registry.Registry, done chan int) {
		assert.JSONEq(t, `{}`, getJSON(t, "/v1/devices/12345"))
		err := reg.Update("12345", device_registry.Device{"D100", -4, 5000, nil})
		assert.NoError(t, err)
		assert.JSONEq(t, `{"instance":"D100", "txPower": -4, "pollPeriod":5000}`, getJSON(t, "/v1/devices/12345"))

		err = reg.Update("12345", device_registry.Device{"D100", -4, 5000, []net.IP{ip}})
		assert.NoError(t, err)
		assert.JSONEq(t, `{"instance":"D100", "txPower": -4, "pollPeriod":5000, "addresses": ["ffff::1"]}`, getJSON(t, "/v1/devices/12345"))
		done <- 1
	})
}

func TestPostV1IP6(t *testing.T) {
	coapServerTest(t, func(t *testing.T, reg *device_registry.Registry, done chan int) {
		dev := device_registry.Device{"D100", -4, 5000, nil}
		err := reg.Update("12345", dev)
		assert.NoError(t, err)

		postJSON(t, "/v1/ip6/12345", `["ffff::1"]`)

		dev2, err := reg.GetOrCreate("12345")
		assert.NoError(t, err)

		dev.Addresses = []net.IP{ip}
		assert.Equal(t, dev, dev2)
		done <- 1
	})
}

func TestGetLastPathPart(t *testing.T) {
	assert.Equal(t, "AABBCCDD", lastPartForPath(t, "/v1/devices/AABBCCDD"))
	assert.Equal(t, "devices", lastPartForPath(t, "/v1/devices/"))
	assert.Equal(t, "devices", lastPartForPath(t, "/v1/devices"))
	assert.Equal(t, "v1", lastPartForPath(t, "/v1"))
	assert.Equal(t, "v1", lastPartForPath(t, "v1"))
}

func coapServerTest(t *testing.T, testFunc func(t *testing.T, reg *device_registry.Registry, done chan int)) {
	reg := device_registry.CreateTestRegistry(t)

	srv, err := NewCoapServer(TEST_COAP_PORT, reg)
	require.NoError(t, err)
	defer srv.Stop()

	testDone := make(chan int, 2)

	go func() {
		err := srv.Serve()
		assert.NoError(t, err)
		testDone <- 1
	}()

	go testFunc(t, reg, testDone)

	<-testDone
}

func getJSON(t *testing.T, path string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := coap_utils.GetJSON(ctx, TEST_COAP_URL, path)
	assert.NoError(t, err)

	return res
}

func postJSON(t *testing.T, path string, payload string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := coap_utils.PostJSON(ctx, TEST_COAP_URL, path, payload)
	assert.Equal(t, "", res)
	assert.NoError(t, err)
}

func lastPartForPath(t *testing.T, path string) string {
	ctx := context.Background()
	poolMsg := pool.AcquireMessage(ctx)
	poolMsg.SetPath(path)
	msg, _ := pool.ConvertTo(poolMsg)

	part, err := coap_utils.GetLastPathPart(&mux.Message{msg, 0})
	assert.NoError(t, err)
	return part
}
