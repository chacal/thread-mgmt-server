package coap_utils

import (
	"context"
	"github.com/pkg/errors"
	"github.com/plgd-dev/go-coap/v2/message"
	"github.com/plgd-dev/go-coap/v2/message/codes"
	"github.com/plgd-dev/go-coap/v2/udp"
	"github.com/plgd-dev/go-coap/v2/udp/client"
	"github.com/plgd-dev/go-coap/v2/udp/message/pool"
	"strings"
	"time"
)

const RequestAckTimeout = 20 * time.Second

func GetJSON(ctx context.Context, url string, path string) (string, error) {
	resp, err := executeRequest(url, path, func() (*pool.Message, error) {
		return client.NewGetRequest(ctx, path)
	})
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

func PostJSON(ctx context.Context, url string, path string, payload string) (string, error) {
	resp, err := executeRequest(url, path, func() (*pool.Message, error) {
		return client.NewPostRequest(ctx, path, message.AppJSON, strings.NewReader(payload))
	})
	if err != nil {
		return "", errors.WithStack(err)
	}

	switch resp.Code() {
	case codes.Empty:
		return "", nil
	case codes.Content:
		body, err := resp.ReadBody()
		if err != nil {
			return "", errors.WithStack(err)
		}
		return string(body), nil
	default:
		return "", errors.Errorf("got response code %v", resp.Code())
	}
}

func executeRequest(url string, path string, reqCreator func() (*pool.Message, error)) (*pool.Message, error) {
	conn, err := udp.Dial(url, udp.WithKeepAlive(nil), udp.WithTransmission(time.Second, RequestAckTimeout, 5))
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't dial to url %v", url)
	}

	req, err := reqCreator()
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't create request with path %v", path)
	}
	req.SetAccept(message.AppJSON)

	resp, err := conn.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return resp, nil
}
