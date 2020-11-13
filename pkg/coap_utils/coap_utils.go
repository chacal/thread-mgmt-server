package coap_utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/plgd-dev/go-coap/v2/message"
	"github.com/plgd-dev/go-coap/v2/message/codes"
	"github.com/plgd-dev/go-coap/v2/mux"
	log "github.com/sirupsen/logrus"
	"strings"
)

func LoggingMiddleware(next mux.Handler) mux.Handler {
	return mux.HandlerFunc(func(w mux.ResponseWriter, r *mux.Message) {
		buf := fmt.Sprintf("%47s | %-19v |", w.Client().RemoteAddr(), r.Code)
		path, err := r.Options.Path()
		if err == nil {
			buf = fmt.Sprintf("%s \"%v\"", buf, path)
		}

		log.Info(buf)
		next.ServeCOAP(w, r)
	})
}

func RespondWithJSON(w mux.ResponseWriter, body interface{}) {
	payload, err := json.Marshal(body)
	if err != nil {
		RespondWithInternalServerError(w, errors.Wrapf(err, "error marshalling payload %v", string(payload)))
		return
	}

	err = setResponse(w, codes.Content, message.AppJSON, payload)
	if err != nil {
		RespondWithInternalServerError(w, errors.WithStack(err))
	}
}

func RespondWithInternalServerError(w mux.ResponseWriter, e error) {
	log.Errorf("%+v", e)
	setResponse(w, codes.InternalServerError, message.TextPlain, nil)
}

func RespondWithBadRequest(w mux.ResponseWriter, e error) {
	log.Errorf("%+v", e)
	setResponse(w, codes.BadRequest, message.TextPlain, nil)
}

func RespondWithEmpty(w mux.ResponseWriter) {
	setResponse(w, codes.Empty, message.TextPlain, nil)
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

func setResponse(w mux.ResponseWriter, code codes.Code, mediaType message.MediaType, payload []byte) error {
	err := w.SetResponse(codes.Content, message.AppJSON, bytes.NewReader(payload))
	if err != nil {
		return errors.WithStack(err)
	}
	logResponse(code, payload)
	return nil
}

func logResponse(code codes.Code, payload []byte) {
	log.Infof("%49s %-19v | %s", "|", code, string(payload))
}