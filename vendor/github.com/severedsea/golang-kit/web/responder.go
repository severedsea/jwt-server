package web

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	"github.com/severedsea/golang-kit/logr"
)

// GenericErrorMessage is the human friendly error message for HTTP 5XX
const GenericErrorMessage = "Sorry, there was a problem. Please try again later."

// RespondJSON writes JSON as http response
func RespondJSON(ctx context.Context, w http.ResponseWriter, object interface{}, headers map[string]string) {
	Responder(ctx, w, object, headers, func(logger logr.Logger, err error) ([]byte, int) {
		resp := NewError(err, err.Error())

		// Handle web.Error
		var werr *Error
		if errors.As(err, &werr) {
			resp = werr
		}

		// Log raw error response
		logger.Errorf("[web/res] Web error: %d %s %s", resp.Status, resp.Code, resp.Desc)

		// 5XX (except 503) should be sanitized before showing to human
		if resp.Status >= 500 && resp.Status != http.StatusServiceUnavailable {
			resp.Desc = GenericErrorMessage
		}

		b, _ := json.Marshal(resp)

		return b, resp.Status
	})
}

type ErrorResponder func(logger logr.Logger, err error) ([]byte, int)

// Responder writes the object/error provided in the argument into the HTTP response
func Responder(ctx context.Context, w http.ResponseWriter, object interface{}, headers map[string]string, errResponder ErrorResponder) {
	logger := logr.GetLogger(ctx)

	// Handle json marshalling error
	respBytes, err := json.Marshal(object)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Errorf("[web/res] JSON marshal error: %s", err)
		return
	}

	// Set HTTP headers
	w.Header().Set("Content-Type", "application/json")
	for key, value := range headers {
		w.Header().Set(key, value)
	}

	// Handle error
	status := http.StatusOK
	if err, ok := object.(error); ok {
		// Add logger fields
		logger = logger.WithField("error", "true")

		respBytes, status = errResponder(logger, err)
	}

	// Log response body
	logger.WithField("status", status).
		Infof("[web/res] %v", string(respBytes))

	// Write response
	w.WriteHeader(status)
	_, _ = w.Write(respBytes)
}
