package util

import (
	"encoding/json"
	"io"
	"net/http"
)

// RequestParseError is an error type indicating that there was a problem parsing the http request
type RequestParseError struct {
	Err     error
	Message string
}

func (e *RequestParseError) Error() string {
	if e.Message != "" {
		return e.Message + " " + e.Err.Error()
	}
	return e.Err.Error()
}

// ParseJsonRequest will parse the request body and return objects in the type of val interface{}
func ParseJsonRequest(r *http.Request, val interface{}) error {
	if r.Body == nil {
		return &RequestParseError{Message: "unable to read message body", Err: nil}
	}
	body, err := io.ReadAll(io.LimitReader(r.Body, 1048576))

	if err != nil {
		return &RequestParseError{Message: "unable to read message body", Err: err}
	}
	defer r.Body.Close()
	err = json.Unmarshal(body, &val)
	if err != nil {
		return &RequestParseError{Message: "unable to parse request", Err: err}
	}
	return nil
}

// GetRequestBody will get the request body as a slice of byte
func GetRequestBody(r *http.Request) ([]byte, error) {

	body, err := io.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		return nil, &RequestParseError{Message: "Unable to read message body", Err: err}
	}
	if err = r.Body.Close(); err != nil {
		return nil, &RequestParseError{Message: "Unable to close message body", Err: err}
	}

	return body, nil
}
