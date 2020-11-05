module github.com/chacal/thread-mgmt-server

go 1.15

require (
	github.com/jessevdk/go-flags v1.4.1-0.20200711081900-c17162fe8fd7
	github.com/pkg/errors v0.9.1
	github.com/plgd-dev/go-coap/v2 v2.1.0
	github.com/sirupsen/logrus v1.4.2
)

replace github.com/plgd-dev/go-coap/v2 v2.1.0 => github.com/chacal/go-coap/v2 v2.1.2-0.20201106113013-5953692068dc
