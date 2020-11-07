package coap_utils

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/plgd-dev/go-coap/v2/message"
	"github.com/plgd-dev/go-coap/v2/message/codes"
	"github.com/plgd-dev/go-coap/v2/mux"
	log "github.com/sirupsen/logrus"
	"strings"
)

func LoggingMiddleware(next mux.Handler) mux.Handler {
	return mux.HandlerFunc(func(w mux.ResponseWriter, r *mux.Message) {
		log.Infof("Client %v, %v", w.Client().RemoteAddr(), r.String())
		next.ServeCOAP(w, r)
	})
}

func RespondWithJSON(w mux.ResponseWriter, body interface{}) {
	payload, err := json.Marshal(body)
	if err != nil {
		RespondWithInternalServerError(w, errors.Wrapf(err, "error marshalling payload %v", string(payload)))
		return
	}

	err = w.SetResponse(codes.Content, message.AppJSON, bytes.NewReader(payload))
	if err != nil {
		RespondWithInternalServerError(w, errors.WithStack(err))
	}
}

func RespondWithInternalServerError(w mux.ResponseWriter, e error) {
	log.Errorf("%+v", e)
	w.SetResponse(codes.InternalServerError, message.TextPlain, nil)
}

func GetLastPathPart(r *mux.Message) (string, error) {
	path, err := r.Message.Options.Path()
	if err != nil {
		return "", errors.Wrapf(err, "couldn't get path from message: %+v", r)
	}

	parts := strings.Split(path, "/")
	if len(parts) < 1 {
		return "", errors.New("empty path")
	}

	return parts[len(parts)-1], nil
}
