package main

import (
	"context"
	"github.com/chacal/thread-mgmt-server/pkg/coap_utils"
	"github.com/chacal/thread-mgmt-server/pkg/device_registry"
	"github.com/chacal/thread-mgmt-server/pkg/test"
	"github.com/plgd-dev/go-coap/v2/mux"
	"github.com/plgd-dev/go-coap/v2/udp/message/pool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestV1Config(t *testing.T) {
	dbFile := test.Tempfile()
	reg, err := device_registry.Open(dbFile)
	require.NoError(t, err)
	defer reg.Close()

	srv, err := NewCoapServer(Options{5683, 8080, dbFile}, reg)
	require.NoError(t, err)
	defer srv.Stop()

	testDone := make(chan int, 2)

	go func() {
		err := srv.Serve()
		assert.NoError(t, err)
		testDone <- 1
	}()

	go func() {
		assert.JSONEq(t, `{}`, getJSON(t, "/v1/config/12345"))
		err = reg.Update("12345", device_registry.Device{"D100", 5000})
		assert.NoError(t, err)
		assert.JSONEq(t, `{"name":"D100", "pollTime":5000}`, getJSON(t, "/v1/config/12345"))
		testDone <- 1
	}()

	<-testDone
}

func TestGetLastPathPart(t *testing.T) {
	assert.Equal(t, "AABBCCDD", lastPartForPath(t, "/v1/config/AABBCCDD"))
	assert.Equal(t, "config", lastPartForPath(t, "/v1/config/"))
	assert.Equal(t, "config", lastPartForPath(t, "/v1/config"))
	assert.Equal(t, "v1", lastPartForPath(t, "/v1"))
	assert.Equal(t, "v1", lastPartForPath(t, "v1"))
}

func getJSON(t *testing.T, path string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := coap_utils.GetJSON(ctx, "localhost:5683", path)
	assert.NoError(t, err)

	return res
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
