package coap_utils

import (
	"context"
	"github.com/pkg/errors"
	"github.com/plgd-dev/go-coap/v2/message"
	"github.com/plgd-dev/go-coap/v2/message/codes"
	"github.com/plgd-dev/go-coap/v2/udp"
	"github.com/plgd-dev/go-coap/v2/udp/client"
)

func GetJSON(ctx context.Context, url string, path string) (string, error) {
	conn, err := udp.Dial(url, udp.WithKeepAlive(nil))
	if err != nil {
		return "", errors.Wrapf(err, "couldn't dial to url %v", url)
	}

	req, err := client.NewGetRequest(ctx, path)
	if err != nil {
		return "", errors.Wrapf(err, "couldn't create request with path %v", path)
	}
	req.SetAccept(message.AppJSON)

	resp, err := conn.Do(req)
	if err != nil {
		return "", errors.WithStack(err)
	}

	if resp.Code() != codes.Content {
		return "", errors.Errorf("got response code %v", resp.Code())
	}

	body, err := resp.ReadBody()
	if err != nil {
		return "", errors.WithStack(err)
	}

	return string(body), nil
}
