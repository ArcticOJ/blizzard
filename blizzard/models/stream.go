package models

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
)

type ResponseStream struct {
	encoder *json.Encoder
	stream  *echo.Response
}

func New(response *echo.Response) *ResponseStream {
	return &ResponseStream{
		encoder: json.NewEncoder(response),
		stream:  response,
	}
}

func (rs *ResponseStream) Write(obj interface{}) error {
	e := rs.encoder.Encode(obj)
	if e != nil {
		return e
	}
	rs.Flush()
	return nil
}

func (rs *ResponseStream) Flush() {
	rs.stream.Flush()
}
