package testutil

import (
	"bytes"
	"io"
	"net/http"
)

func SomeRes(body []byte) *http.Response {
	return &http.Response{
		Status:        "200 OK",
		StatusCode:    200,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Body:          io.NopCloser(bytes.NewBuffer(body)),
		ContentLength: int64(len(body)),
		Header:        make(http.Header, 0),
	}
}
