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
	"strings"
	"testing"
)

var LOCAL_IP = gonet.ParseIP("127.0.0.1")
var ip = gonet.ParseIP("ffff::1")
var testState = device_registry.State{[]gonet.IP{ip}, 2970, "A100",
	device_registry.ParentInfo{"0x4400", 3, 2, -65, -63},
}

func TestGateway_PushSettings(t *testing.T) {
	testWithCoapServer(t, func(t *testing.T, r *mux.Router, done chan int) {
		expectJSONPost(t, r, "api/settings", `{"instance": "D100","txPower": -4,"pollPeriod": 5000}`)

		gw := Create()
		dev := device_registry.Defaults{"D100", -4, 5000}
		err := gw.PushDefaults(dev, LOCAL_IP)
		assert.NoError(t, err)
		done <- 1
	})
}

func TestGateway_FetchState(t *testing.T) {
	testWithCoapServer(t, func(t *testing.T, r *mux.Router, done chan int) {
		expectJSONGet(t, r, "api/state",
			`{
				"vcc": 2970,
				"instance": "A100",
				"addresses": [
					"ffff::1"
				],
				"parent": {
					"rloc16": "0x4400",
					"linkQualityIn": 3,
					"linkQualityOut": 2,
					"avgRssi": -65,
					"latestRssi": -63
				}
			}`,
		)

		gw := Create()
		state, err := gw.FetchState(LOCAL_IP)
		assert.NoError(t, err)
		assert.Equal(t, testState, state)
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

func expectJSONPost(t *testing.T, r *mux.Router, path string, body string) {
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

func expectJSONGet(t *testing.T, r *mux.Router, path string, response string) {
	_ = r.Handle(path, mux.HandlerFunc(func(w mux.ResponseWriter, msg *mux.Message) {
		assert.Equal(t, codes.GET, msg.Code)
		p, _ := msg.Options.Path()
		assert.Equal(t, path, p)
		cf, _ := msg.Options.Accept()
		assert.Equal(t, message.AppJSON, cf)
		_ = w.SetResponse(codes.Content, message.AppJSON, strings.NewReader(response))
	}))
}
