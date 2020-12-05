package device_gateway

import (
	"github.com/chacal/thread-mgmt-server/pkg/device_registry"
	"github.com/plgd-dev/go-coap/v2/message"
	"github.com/plgd-dev/go-coap/v2/message/codes"
	"github.com/plgd-dev/go-coap/v2/mux"
	"github.com/plgd-dev/go-coap/v2/net"
	"github.com/plgd-dev/go-coap/v2/udp"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	gonet "net"
	"testing"
)

var LOCAL_IP = gonet.ParseIP("127.0.0.1")

func TestGateway_PushSettings(t *testing.T) {
	testWithCoapServer(t, func(t *testing.T, r *mux.Router, done chan int) {
		assertJSONPost(t, r, "api/settings", `{"instance": "D100","txPower": -4,"pollPeriod": 5000}`)

		gw := Create()
		dev := device_registry.Device{"D100", -4, 5000, []device_registry.DeviceAddress{{LOCAL_IP, false}}}
		err := gw.PushSettings(dev, LOCAL_IP)
		assert.NoError(t, err)
		done <- 1
	})
}

func testWithCoapServer(t *testing.T, testFunc func(t *testing.T, r *mux.Router, done chan int)) {
	r := mux.NewRouter()

	srv := udp.NewServer(udp.WithMux(r), udp.WithKeepAlive(nil))
	defer srv.Stop()

	conn, err := net.NewListenUDP("udp", ":5683")
	assert.NoError(t, err)
	defer conn.Close()

	testDone := make(chan int, 2)

	go func() {
		err := srv.Serve(conn)
		assert.NoError(t, err)
		testDone <- 1
	}()

	go testFunc(t, r, testDone)

	<-testDone
}

func assertJSONPost(t *testing.T, r *mux.Router, path string, body string) {
	_ = r.Handle(path, mux.HandlerFunc(func(w mux.ResponseWriter, msg *mux.Message) {
		assert.Equal(t, codes.POST, msg.Code)
		p, _ := msg.Options.Path()
		assert.Equal(t, path, p)
		cf, _ := msg.Options.ContentFormat()
		assert.Equal(t, message.AppJSON, cf)
		b, _ := ioutil.ReadAll(msg.Body)
		assert.JSONEq(t, body, string(b))
		_ = w.SetResponse(codes.Empty, message.TextPlain, nil)
	}))
}
