package httputils

import "errors"

const (
	statusFailedHttpRequest  = "http request error:"
	statusFailedHttpResponse = "http read response error:"
	statusTimeout            = "http request timeout:"
)

const (
	MIMEApplicationJSON = "application/json"
	ContentTypeHeader   = "Content-Type"
)

var (
	ErrEmptyResponse = errors.New("empty response")
)
