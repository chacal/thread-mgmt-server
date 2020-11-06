package main

import (
	"context"
	"fmt"
	"github.com/plgd-dev/go-coap/v2/message/codes"
	"github.com/plgd-dev/go-coap/v2/net"
	"github.com/plgd-dev/go-coap/v2/udp"
	"github.com/plgd-dev/go-coap/v2/udp/client"
	"github.com/plgd-dev/go-coap/v2/udp/message/pool"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var testOpts = Options{"en0", "[ff03::1]", 9999, "test.mgmt"}
var expectedResp = fmt.Sprintf("{\"mgmtServer\": \"%v\"}", testOpts.MgmtServerAddress)

func TestDiscovery(t *testing.T) {
	log.SetLevel(log.WarnLevel)
	done := make(chan int)

	go startServer(t, done)

	go func() {
		url := fmt.Sprint(testOpts.ListenAddr, ":", testOpts.Port)
		startDiscoveryClient(t, url, "/discover", func(cc *client.ClientConn, resp *pool.Message) {
			b, _ := resp.Message.ReadBody()
			assert.Equal(t, codes.Content, resp.Code())
			assert.JSONEq(t, expectedResp, string(b))
			done <- 1
		})
	}()

	<-done
}

func startServer(t *testing.T, done chan int) {
	defer close(done)
	err := startCoapServer(testOpts)
	require.NoError(t, err)
}

func startDiscoveryClient(t *testing.T, url string, path string, responseHandler func(cc *client.ClientConn, resp *pool.Message)) {
	conn, err := net.NewListenUDP("udp6", "")
	require.NoError(t, err)
	defer conn.Close()

	server := udp.NewServer()
	defer server.Stop()

	go func() {
		err := server.Serve(conn)
		require.NoError(t, err)
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = server.Discover(ctx, url, path, responseHandler)
	require.NoError(t, err)
}