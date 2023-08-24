package web

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
	"github.com/severedsea/golang-kit/logr"
	"github.com/severedsea/golang-kit/redact"
)

// ParseJSONBody parses JSON request body and handles errors
func ParseJSONBody(result interface{}, r io.ReadCloser) ([]byte, error) {
	return ParseAndLogJSONBody(nil, result, r, redact.Keys{})
}

// ParseAndLogJSONBody parses JSON request body and handles errors
func ParseAndLogJSONBody(logger logr.Logger, result interface{}, r io.ReadCloser, redactKeys redact.Keys) ([]byte, error) {
	reqBytes, err := ioutil.ReadAll(r)
	defer r.Close()
	if err != nil {
		return nil, &Error{Status: http.StatusBadRequest, Code: "read_body", Desc: err.Error(), Err: errors.WithStack(err)}
	}

	if logger != nil {
		bytes := reqBytes

		if !redactKeys.IsEmpty() {
			// Redact request body before logging
			var jsonBody map[string]interface{}
			if err := json.Unmarshal(bytes, &jsonBody); err != nil {
				return nil, &Error{Status: http.StatusBadRequest, Code: "parse_body", Desc: err.Error(), Err: errors.WithStack(err)}
			}

			redactedJSONBody := redact.MaskMap(jsonBody, redactKeys)
			bytes, err = json.Marshal(redactedJSONBody)
			if err != nil {
				return nil, &Error{Status: http.StatusBadRequest, Code: "parse_body", Desc: err.Error(), Err: errors.WithStack(err)}
			}
		}

		logger.Infof("[web/req] %s", string(bytes))
	}

	err = json.Unmarshal(reqBytes, &result)
	if err != nil {
		return nil, &Error{Status: http.StatusBadRequest, Code: "parse_body", Desc: err.Error(), Err: errors.WithStack(err)}
	}

	return reqBytes, nil
}
