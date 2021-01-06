package mqtt

import (
	"fmt"
	"github.com/chacal/thread-mgmt-server/pkg/device_registry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net"
	"testing"
	"time"
)

var ip = net.ParseIP("ffff::1")
var testState = device_registry.State{[]net.IP{ip}, 2970, "A100", -4, 1000,
	device_registry.ParentInfo{"0x4400", 3, 2, -65, -63},
}

func TestMqttSender_publishDataForState(t *testing.T) {
	ts := time.Now()
	formattedTs := ts.Format(time.RFC3339)

	topic, payload, err := publishDataForState(testState, ts)
	require.NoError(t, err)
	assert.Equal(t, "/sensor/A100/d/state", topic)
	assert.JSONEq(t,
		fmt.Sprintf(`{
			"instance": "A100",
			"tag": "d",
			"ts": "%v",
			"vcc": 2970,
			"parent": {
				"rloc16": "0x4400",
				"linkQualityIn": 3,
				"linkQualityOut": 2,
				"avgRssi": -65,
				"latestRssi": -63
			}
		}`, formattedTs),
		string(payload),
	)
}
